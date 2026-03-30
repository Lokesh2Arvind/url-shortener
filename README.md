# 🔗 URL Shortener (Go + Redis)

A high-performance URL shortener built using **Go** and **Redis**, supporting custom aliases, expiration, and real-time analytics.

---

## 🚀 Features

  - **Shorten long URLs** into compact, shareable links
  - **Custom alias support** (e.g., `/github`, `/docs`)
  - **Automatic expiration** using Redis TTL
 - **Click analytics** (track number of visits per link)
 - **Update alias** without breaking existing mappings
 - **Delete URLs** with associated metadata
 - **Automatic routing** (`/abc123` → redirect)
 - **Atomic ID generation** using Redis (no collisions)

---

## 🧠 How It Works

1. User submits a long URL
2. Server:

   * Generates a unique short code (or uses custom alias)
   * Stores mapping in Redis
3. On visiting the short URL:

   * Server fetches original URL from Redis
   * Redirects user
   * Increments click counter

---

## ⚙️ Tech Stack

* **Backend:** Go (Golang)
* **Database:** Redis
* **Frontend:** HTML + CSS
* **Architecture:** REST-style handlers with key-value storage

---

## 🗄️ Data Model (Redis)

```
short:<code>   → original URL
long:<url>     → short code
clicks:<code>  → number of visits
counter        → global ID generator
```

---

## ✨ Unique Aspects

* ⚡ **Atomic ID generation using Redis `INCR`**
  Ensures globally unique short URLs without collisions

* ⏱ **Per-link expiration using Redis TTL**
  Links automatically expire without manual cleanup

* 🔁 **Bidirectional mapping (short ↔ long)**
  Prevents duplicate URL shortening

* 📊 **Real-time analytics tracking**
  Lightweight click counting using Redis

---

## 📦 Setup & Run

### 1. Clone the repository

```bash
git clone https://github.com/your-username/url-shortener.git
cd url-shortener
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Start Redis

```bash
redis-server
```

### 4. Run the server

```bash
go run main.go
```

### 5. Open in browser

```
http://localhost:8080
```

---

## 📌 Example

```
Input:  https://google.com
Output: http://localhost:8080/abc123
```

---

## 📈 Future Improvements

 - Rate limiting using Redis
 - User authentication
 - Dashboard UI for analytics
 - QR code generation

---
## Project URL 

https://roadmap.sh/projects/url-shortening-service

---

## 👨‍💻 Author

**Lokesh Arvind**

---

## ⭐ If you like this project

Give it a star on GitHub!
