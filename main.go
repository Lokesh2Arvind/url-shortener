package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"urlshort/shorten"

	"github.com/redis/go-redis/v9"
)

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

/*
Key Design:

short:<code> → long URL
long:<url> → short code
clicks:<code> → count
*/

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	url := r.FormValue("url")
	alias := r.FormValue("alias")
	expiryStr := r.FormValue("expiry")

	// Validate expiry
	expiry := 0
	if expiryStr != "" {
		val, err := strconv.Atoi(expiryStr)
		if err != nil {
			http.Error(w, "Invalid expiry value", http.StatusBadRequest)
			return
		}
		expiry = val
	}

	// Check if URL already shortened
	existing, err := rdb.Get(ctx, "long:"+url).Result()
	if err == nil {
		http.Error(w, "Already exists: http://localhost:8080/"+existing, http.StatusConflict)
		return
	}

	var finalName string

	// Handle alias
	if alias != "" {
		_, err := rdb.Get(ctx, "short:"+alias).Result()
		if err == nil {
			http.Error(w, "Alias already taken", http.StatusConflict)
			return
		}
		finalName = alias
	}

	// Generate ID if no alias
	if finalName == "" {
		id, _ := rdb.Incr(ctx, "counter").Result()
		finalName = shorten.Shorten(url, id)
	}

	// TTL handling
	var ttl time.Duration = 0
	if expiry > 0 {
		ttl = time.Duration(expiry) * 24 * time.Hour
	}

	// Store mappings
	rdb.Set(ctx, "short:"+finalName, url, ttl)
	rdb.Set(ctx, "long:"+url, finalName, ttl)

	fmt.Fprintf(w, `
	<html>
	<body>
		<h3>Short URL:</h3>
		<a href="/%s">http://localhost:8080/%s</a>
		<br><br>
		<a href="/">Go Back</a>
	</body>
	</html>`,
		finalName, finalName)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	short := r.URL.Path[1:]

	long, err := rdb.Get(ctx, "short:"+short).Result()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// increment clicks
	rdb.Incr(ctx, "clicks:"+short)

	http.Redirect(w, r, long, http.StatusFound)
}

func changeAlias(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	oldAlias := r.FormValue("old_alias")
	newAlias := r.FormValue("new_alias")

	longURL, err := rdb.Get(ctx, "short:"+oldAlias).Result()
	if err != nil {
		http.Error(w, "Old alias not found", http.StatusNotFound)
		return
	}

	_, err = rdb.Get(ctx, "short:"+newAlias).Result()
	if err == nil {
		http.Error(w, "New alias already exists", http.StatusConflict)
		return
	}

	// move data
	rdb.Set(ctx, "short:"+newAlias, longURL, 0)
	rdb.Del(ctx, "short:"+oldAlias)

	// update reverse mapping
	rdb.Set(ctx, "long:"+longURL, newAlias, 0)

	// move clicks
	clicks, _ := rdb.Get(ctx, "clicks:"+oldAlias).Result()
	if clicks != "" {
		rdb.Set(ctx, "clicks:"+newAlias, clicks, 0)
		rdb.Del(ctx, "clicks:"+oldAlias)
	}

	fmt.Fprintf(w, `
	<html>
	<body>
		<h3>Alias Updated</h3>
		<p>%s → %s</p>
		<a href="/">Go Back</a>
	</body>
	</html>`,
		oldAlias, newAlias)
}

func removeAlias(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	alias := r.FormValue("alias")

	longURL, err := rdb.Get(ctx, "short:"+alias).Result()
	if err != nil {
		http.Error(w, "Alias not found", http.StatusNotFound)
		return
	}

	clicks, _ := rdb.Get(ctx, "clicks:"+alias).Result()

	rdb.Del(ctx, "short:"+alias)
	rdb.Del(ctx, "long:"+longURL)
	rdb.Del(ctx, "clicks:"+alias)

	fmt.Fprintf(w, `
	<html>
	<body>
		<h3>Deleted</h3>
		<p>%s → %s (Clicks: %s)</p>
		<a href="/">Go Back</a>
	</body>
	</html>`,
		alias, longURL, clicks)
}

func allHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	html := `<html><body><h2>All URLs</h2><ul>`

	keys, _ := rdb.Keys(ctx, "short:*").Result()

	for _, key := range keys {

		short := key[6:] // remove "short:"
		long, _ := rdb.Get(ctx, key).Result()
		clicks, _ := rdb.Get(ctx, "clicks:"+short).Result()

		html += fmt.Sprintf(
			`<li><a href="/%s">%s</a> → %s (Clicks: %s)</li>`,
			short, short, long, clicks,
		)
	}

	html += `</ul></body></html>`

	fmt.Fprintln(w, html)
}

func main() {

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/update", changeAlias)
	http.HandleFunc("/delete", removeAlias)
	http.HandleFunc("/all", allHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./static/index.html")
			return
		}
		redirectHandler(w, r)
	})

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}
