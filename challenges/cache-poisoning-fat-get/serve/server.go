package serve

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type cache struct {
	mu    sync.Mutex
	store map[string]string
}

func newCache() *cache {
	return &cache{store: make(map[string]string)}
}

func RunServer(port string, vulnerable bool) {
	c := newCache()

	// /search is a "fat GET": it accepts the query term both as a query
	// string parameter and, non-standardly, as a request body. The cache
	// keys purely on the URL (as caches normally do, since GET bodies are
	// not expected to affect the response), but the vulnerable origin still
	// reads the body if present and lets it override the query parameter.
	//
	// An attacker can therefore send a GET with a body to a URL that's
	// already cached (or about to be) and have the cache store the
	// body-influenced response under the body-less cache key, poisoning it
	// for every subsequent caller of that same URL.
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path + "?" + r.URL.RawQuery

		c.mu.Lock()
		body, hit := c.store[key]
		c.mu.Unlock()

		w.Header().Set("Content-Type", "text/html")
		if hit {
			w.Header().Set("X-Cache", "HIT")
			w.Write([]byte(body))
			return
		}

		term := r.URL.Query().Get("q")
		if vulnerable {
			if b, err := io.ReadAll(r.Body); err == nil && len(b) > 0 {
				term = string(b)
			}
		}

		body = fmt.Sprintf("<html><body>Results for: %s</body></html>", term)

		c.mu.Lock()
		c.store[key] = body
		c.mu.Unlock()

		w.Header().Set("X-Cache", "MISS")
		w.Write([]byte(body))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
