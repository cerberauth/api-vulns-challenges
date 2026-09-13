package serve

import (
	"log"
	"net/http"
	"strings"
	"sync"
)

type cacheEntry struct {
	body string
}

type cache struct {
	mu    sync.Mutex
	store map[string]cacheEntry
}

func newCache() *cache {
	return &cache{store: make(map[string]cacheEntry)}
}

func RunServer(port string, vulnerable bool) {
	c := newCache()

	// /home renders content whose language depends on the caller's session
	// cookie, but the cache key (and the Vary declaration) do not account
	// for it. The first cookie-bearing caller's personalized response gets
	// stored under the shared, cookie-agnostic key and served to every
	// other visitor regardless of their own session.
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		cookie := r.Header.Get("Cookie")

		key := r.URL.Path
		if !vulnerable {
			key += "|" + cookie
		}

		c.mu.Lock()
		entry, hit := c.store[key]
		c.mu.Unlock()

		w.Header().Set("Content-Type", "text/html")
		if vulnerable {
			// Vary is absent even though the response demonstrably differs
			// by Cookie - the precondition RFC 9111 requires callers (and
			// caches) be told about via Vary
		} else {
			w.Header().Set("Vary", "Cookie")
		}

		if hit {
			w.Header().Set("X-Cache", "HIT")
			w.Write([]byte(entry.body))
			return
		}

		greeting := "Welcome, guest"
		if strings.Contains(cookie, "session=") {
			name := strings.TrimPrefix(cookie, "session=")
			greeting = "Welcome back, " + name
		}
		body := "<html><body>" + greeting + "</body></html>"

		c.mu.Lock()
		c.store[key] = cacheEntry{body: body}
		c.mu.Unlock()

		w.Header().Set("X-Cache", "MISS")
		w.Write([]byte(body))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
