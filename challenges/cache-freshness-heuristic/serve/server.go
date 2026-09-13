package serve

import (
	"log"
	"net/http"
	"sync"
)

type cache struct {
	mu    sync.Mutex
	store map[string]cachedResponse
}

type cachedResponse struct {
	status int
	body   string
}

func newCache() *cache {
	return &cache{store: make(map[string]cachedResponse)}
}

func RunServer(port string, vulnerable bool) {
	c := newCache()

	// /transient always errors with a 503 and carries no explicit caching
	// directives. RFC 9110 §15.1 only lists 200, 203, 204, 206, 300, 301,
	// 404, 405, 410, 414 and 501 as heuristically cacheable by default; 503
	// is not among them. The vulnerable cache heuristically caches it
	// anyway, so a single transient origin failure gets replayed to every
	// caller long after the origin has recovered.
	http.HandleFunc("/transient", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path

		c.mu.Lock()
		resp, hit := c.store[key]
		c.mu.Unlock()
		if hit {
			w.Header().Set("X-Cache", "HIT")
			w.WriteHeader(resp.status)
			w.Write([]byte(resp.body))
			return
		}

		status := http.StatusServiceUnavailable
		body := "503 Service Unavailable (transient)"

		if vulnerable {
			// heuristically cached even though the status code and the
			// absence of any Cache-Control/Expires header do not permit it
			c.mu.Lock()
			c.store[key] = cachedResponse{status: status, body: body}
			c.mu.Unlock()
		}

		w.Header().Set("X-Cache", "MISS")
		w.WriteHeader(status)
		w.Write([]byte(body))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
