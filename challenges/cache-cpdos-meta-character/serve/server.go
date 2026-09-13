package serve

import (
	"log"
	"net/http"
	"strings"
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

// metaCharacters are control characters a CDN happily forwards after
// percent-decoding a query value, but that the origin's stricter input
// validation rejects outright.
var metaCharacters = []string{"\r", "\n", "\x00", "\a"}

func containsMetaCharacter(v string) bool {
	for _, ch := range metaCharacters {
		if strings.Contains(v, ch) {
			return true
		}
	}
	return false
}

func RunServer(port string, vulnerable bool) {
	c := newCache()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// the cache keys purely on the path, unaware of the query value's
		// content, while the origin below actually inspects it
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

		name := r.URL.Query().Get("name")

		var status int
		var body string
		if containsMetaCharacter(name) {
			// the origin's stricter input validation rejects control
			// characters in this field outright
			status = http.StatusBadRequest
			body = "400 Bad Request: invalid character in name"
		} else {
			status = http.StatusOK
			body = "welcome, " + name
		}

		if vulnerable && status != http.StatusOK {
			// the CDN cached the resulting origin error keyed on the path
			// alone, unaware the query value that triggered it even
			// existed - CPDoS "HTTP Meta Character" (HMC)
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
