package serve

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

type cacheEntry struct {
	body        string
	contentType string
}

// cache simulates a shared edge cache in front of the backend.
type cache struct {
	mu    sync.Mutex
	store map[string]cacheEntry
}

func newCache() *cache {
	return &cache{store: make(map[string]cacheEntry)}
}

func RunServer(port string, vulnerable bool) {
	c := newCache()

	// /account returns data scoped to the caller's Authorization token.
	http.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		// vulnerable: the cache key only considers the path, ignoring the
		// Authorization header, so the first caller's authenticated response
		// gets served to every subsequent caller regardless of their token
		key := r.URL.Path
		if !vulnerable {
			key = r.URL.Path + "|" + token
		}

		c.mu.Lock()
		entry, hit := c.store[key]
		c.mu.Unlock()

		if hit {
			w.Header().Set("Content-Type", entry.contentType)
			w.Header().Set("X-Cache", "HIT")
			w.Write([]byte(entry.body))
			return
		}

		body := fmt.Sprintf(`{"email": "%s@example.com"}`, strings.TrimPrefix(token, "Bearer "))
		entry = cacheEntry{body: body, contentType: "application/json"}

		c.mu.Lock()
		c.store[key] = entry
		c.mu.Unlock()

		w.Header().Set("Content-Type", entry.contentType)
		w.Header().Set("X-Cache", "MISS")
		w.Write([]byte(entry.body))
	})

	// /assets/profile.js simulates web cache deception: an authenticated
	// endpoint's content is exposed through a path that looks like a static
	// asset, which shared caches will happily cache and serve to anyone.
	http.HandleFunc("/assets/profile.js", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/javascript")

		if vulnerable {
			// cached purely by path, so the response containing the first
			// authenticated caller's data is reused for anonymous requests
			key := r.URL.Path
			c.mu.Lock()
			entry, hit := c.store[key]
			c.mu.Unlock()
			if hit {
				w.Header().Set("X-Cache", "HIT")
				w.Write([]byte(entry.body))
				return
			}
			body := fmt.Sprintf("var profile = {token: %q};", token)
			c.mu.Lock()
			c.store[key] = cacheEntry{body: body}
			c.mu.Unlock()
			w.Header().Set("X-Cache", "MISS")
			w.Write([]byte(body))
			return
		}

		// fixed: authenticated content is never cached, and is not served
		// under a static-looking path in the first place
		w.Header().Set("Cache-Control", "no-store")
		http.NotFound(w, r)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
