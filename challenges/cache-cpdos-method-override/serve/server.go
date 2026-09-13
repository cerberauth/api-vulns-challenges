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

// methodOverrideHeaders are honored by the vulnerable origin to let a GET
// request be re-interpreted as a different, possibly blocked, method - a
// pattern some frameworks support for HTML forms that can't send PUT/DELETE
// directly.
var methodOverrideHeaders = []string{"X-HTTP-Method-Override", "X-Method-Override", "X-HTTP-Method"}

// blockedMethods are rejected by the origin with an error, per its routing
// rules (only GET is exposed on this endpoint).
var blockedMethods = map[string]bool{"DELETE": true, "PUT": true, "PATCH": true}

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

		effectiveMethod := r.Method
		if vulnerable {
			for _, h := range methodOverrideHeaders {
				if v := r.Header.Get(h); v != "" {
					effectiveMethod = v
					break
				}
			}
		}

		var status int
		var body string
		if blockedMethods[effectiveMethod] {
			// the origin rejects the (overridden) method with an error
			status = http.StatusMethodNotAllowed
			body = "405 Method Not Allowed"
		} else {
			status = http.StatusOK
			body = "welcome"
		}

		if vulnerable && status != http.StatusOK {
			// the CDN caches keyed on GET /, unaware that a method-override
			// header changed the effective method the origin acted on -
			// CPDoS "HTTP Method Override" (HMO)
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
