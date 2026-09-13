package serve

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

// cache simulates a shared edge cache keyed only by path (as most CDNs do
// by default), fronting the origin below.
type cache struct {
	mu    sync.Mutex
	store map[string]string
}

func newCache() *cache {
	return &cache{store: make(map[string]string)}
}

func RunServer(port string, vulnerable bool) {
	c := newCache()

	// / renders a page that links back to itself using a "canonical" host.
	// The origin builds that absolute link from the caller-supplied
	// X-Forwarded-Host header instead of a fixed, trusted config value.
	//
	// Because the edge cache keys only on the path, whichever caller's
	// X-Forwarded-Host value is used to build the very first cached
	// response gets served to every subsequent visitor of "/" until the
	// entry expires - a classic unkeyed-input web cache poisoning
	// primitive (Kettle, "Practical Web Cache Poisoning", 2018).
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path

		c.mu.Lock()
		body, hit := c.store[key]
		c.mu.Unlock()

		w.Header().Set("Content-Type", "text/html")
		if hit {
			w.Header().Set("X-Cache", "HIT")
			w.Write([]byte(body))
			return
		}

		host := "example.com"
		if vulnerable {
			if xfh := r.Header.Get("X-Forwarded-Host"); xfh != "" {
				host = xfh
			}
		}
		body = fmt.Sprintf(`<html><head><link rel="canonical" href="https://%s/"></head>`+
			`<body><script src="https://%s/static/app.js"></script></body></html>`, host, host)

		if vulnerable {
			c.mu.Lock()
			c.store[key] = body
			c.mu.Unlock()
		}

		w.Header().Set("X-Cache", "MISS")
		w.Write([]byte(body))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
