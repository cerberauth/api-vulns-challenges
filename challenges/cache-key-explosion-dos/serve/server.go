package serve

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

// maxDistinctKeys caps how many distinct cache entries the (fixed) cache
// allows for a given path, beyond which it falls back to a single, shared
// entry rather than growing without bound.
const maxDistinctKeys = 16

type cache struct {
	mu    sync.Mutex
	store map[string]int
}

func newCache() *cache {
	return &cache{store: make(map[string]int)}
}

func RunServer(port string, vulnerable bool) {
	c := newCache()
	var originHits int
	var originHitsMu sync.Mutex

	// /search takes an arbitrary "q" query parameter, and, in the
	// vulnerable configuration, the cache keys on the full, unbounded query
	// string - so an attacker can generate an unlimited number of distinct
	// cache entries (or, if the cache is small, force an eviction storm and
	// an origin request for every single value) simply by varying an
	// unkeyed-in-practice parameter across many requests.
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")

		key := q
		if !vulnerable {
			// fixed: cap the number of distinct cache entries kept per
			// path; once the limit is reached, further distinct values
			// share a single "overflow" bucket instead of exploding the
			// cache's storage or hammering the origin per unique value
			c.mu.Lock()
			if _, exists := c.store[key]; !exists && len(c.store) >= maxDistinctKeys {
				key = "__overflow__"
			}
			c.mu.Unlock()
		}

		c.mu.Lock()
		_, hit := c.store[key]
		if !hit {
			c.store[key] = 0
		}
		c.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if hit {
			w.Header().Set("X-Cache", "HIT")
		} else {
			w.Header().Set("X-Cache", "MISS")
			originHitsMu.Lock()
			originHits++
			originHitsMu.Unlock()
		}

		w.Header().Set("X-Distinct-Cache-Entries", fmt.Sprintf("%d", len(c.store)))
		originHitsMu.Lock()
		hits := originHits
		originHitsMu.Unlock()
		w.Write([]byte(fmt.Sprintf(`{"results": [], "origin_requests_so_far": %d}`, hits)))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
