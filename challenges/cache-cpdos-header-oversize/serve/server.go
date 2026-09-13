package serve

import (
	"log"
	"net/http"
	"sync"
)

// originHeaderLimit mimics a typical origin's header-block size limit
// (e.g. Apache's default of 8192 bytes). The CDN in front of it, in the
// vulnerable configuration, accepts much larger header blocks (e.g.
// CloudFront's 20480 bytes) and forwards them straight through - a size
// mismatch that lets an attacker craft a request the CDN accepts but the
// origin rejects with an error.
const originHeaderLimit = 8192

// cdnHeaderLimit is the shared cache/CDN's own, larger limit.
const cdnHeaderLimit = 20480

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

func headerBlockSize(r *http.Request) int {
	size := 0
	for name, values := range r.Header {
		for _, v := range values {
			size += len(name) + len(v)
		}
	}
	return size
}

func RunServer(port string, vulnerable bool) {
	c := newCache()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

		size := headerBlockSize(r)

		// the CDN enforces its own, larger limit before ever forwarding
		// the request upstream
		if size > cdnHeaderLimit {
			http.Error(w, "request header fields too large", http.StatusRequestHeaderFieldsTooLarge)
			return
		}

		var status int
		var body string
		if size > originHeaderLimit {
			// the origin rejects a header block this large
			status = http.StatusRequestHeaderFieldsTooLarge
			body = "431 Request Header Fields Too Large"
		} else {
			status = http.StatusOK
			body = "welcome"
		}

		if vulnerable && status != http.StatusOK {
			// the CDN caches this error response and replays it to every
			// subsequent, legitimate caller of the same URL - a denial of
			// service via HTTP Header Oversize (CPDoS "HHO")
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
