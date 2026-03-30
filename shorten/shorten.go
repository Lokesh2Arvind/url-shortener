package shorten

import (
	"log"
	"net/url"
	"strings"
)

func base62Encode(num int64) string {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var result string = ""
	for num > 0 {
		remainder := num % 62
		result = string(chars[remainder]) + result
		num /= 62
	}
	return result
}

// Shorten generates a short string based on the input URL. Its an exported function
func Shorten(rurl string, id int64) string {
	parsed, _ := url.Parse(rurl)

	host := parsed.Host
	arr := strings.Split(host, ".")
	var name string
	if arr[0] == "www" {
		name = arr[1]
	} else {
		name = arr[0]
	}
	random := base62Encode(id)
	log.Printf("Generated short URL: %s for long URL: %s", name+random, rurl)
	return name + random
}
