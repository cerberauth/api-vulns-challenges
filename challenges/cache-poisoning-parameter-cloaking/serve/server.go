package serve

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

type cache struct {
	mu    sync.Mutex
	store map[string]string
}

func newCache() *cache {
	return &cache{store: make(map[string]string)}
}

// parseExcludedVulnerable mimics a shared cache that keys on the "id" query
// parameter but treats everything after a ";" as opaque path/matrix
// parameters (Rails-style), so it never looks past the ";" when building the
// cache key.
func cacheKeyVulnerable(r *http.Request) string {
	raw := r.URL.RawQuery
	if idx := strings.Index(raw, ";"); idx != -1 {
		raw = raw[:idx]
	}
	return r.URL.Path + "?" + raw
}

// cacheKeyFixed keys on the full, decoded query string so a value smuggled
// in after a ";" is still part of the key.
func cacheKeyFixed(r *http.Request) string {
	return r.URL.Path + "?" + r.URL.RawQuery
}

func RunServer(port string, vulnerable bool) {
	c := newCache()

	// /greet reflects the "name" query parameter into the response. The
	// application itself (unlike the cache) parses ";" as an additional
	// parameter separator, so "?id=1;name=attacker" is read by the app as
	// name=attacker while the cache only sees "id=1" as the key - letting an
	// attacker "cloak" an unkeyed payload behind a keyed-looking parameter
	// and poison the shared cache entry for id=1.
	//
	// Kettle, "Web Cache Entanglement" (Black Hat USA 2020).
	http.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		var key string
		if vulnerable {
			key = cacheKeyVulnerable(r)
		} else {
			key = cacheKeyFixed(r)
		}

		c.mu.Lock()
		body, hit := c.store[key]
		c.mu.Unlock()

		w.Header().Set("Content-Type", "text/html")
		if hit {
			w.Header().Set("X-Cache", "HIT")
			w.Write([]byte(body))
			return
		}

		name := "world"
		if vulnerable {
			// the application splits on ";" the same way the origin's
			// framework would, so a value cloaked after the cache's
			// delimiter is still honored
			raw := r.URL.RawQuery
			for _, part := range strings.Split(raw, ";") {
				if strings.HasPrefix(part, "name=") {
					name = strings.TrimPrefix(part, "name=")
				}
			}
		} else if v := r.URL.Query().Get("name"); v != "" {
			name = v
		}

		body = fmt.Sprintf("<html><body>Hello, %s!</body></html>", name)

		c.mu.Lock()
		c.store[key] = body
		c.mu.Unlock()

		w.Header().Set("X-Cache", "MISS")
		w.Write([]byte(body))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
