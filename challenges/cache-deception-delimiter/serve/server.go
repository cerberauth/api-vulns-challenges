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

// delimiters the shared cache treats as "the origin must have stopped
// parsing the path here, so anything after this is a static suffix" - while
// the origin itself ignores them entirely and keeps routing on the full
// path up to the delimiter.
var originOnlyDelimiters = []string{";", "%00", "%0a", "$"}

func RunServer(port string, vulnerable bool) {
	c := newCache()

	// registered directly (bypassing http.ServeMux) so delimiter
	// characters in the raw path reach the handler unmodified
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawPath := r.URL.Path

		if !strings.HasPrefix(rawPath, "/account") {
			http.NotFound(w, r)
			return
		}

		cacheable := false
		if vulnerable {
			for _, d := range originOnlyDelimiters {
				if idx := strings.Index(rawPath, d); idx != -1 {
					suffix := rawPath[idx+len(d):]
					// the cache sees a static-looking suffix after the
					// delimiter (e.g. ".css") and decides to cache the
					// whole response, even though the delimiter and
					// suffix carry no meaning to the origin
					if strings.HasSuffix(suffix, ".css") || strings.HasSuffix(suffix, ".js") {
						cacheable = true
					}
					break
				}
			}
		}

		token := r.Header.Get("Authorization")
		body := fmt.Sprintf(`{"email": "%s@example.com", "session_token": "%s"}`, strings.TrimPrefix(token, "Bearer "), token)

		if cacheable {
			// the origin ignores everything from the delimiter onward and
			// serves the same authenticated /account response regardless
			key := rawPath
			c.mu.Lock()
			cached, hit := c.store[key]
			c.mu.Unlock()
			if hit {
				w.Header().Set("X-Cache", "HIT")
				w.Write([]byte(cached))
				return
			}
			c.mu.Lock()
			c.store[key] = body
			c.mu.Unlock()
			w.Header().Set("X-Cache", "MISS")
			w.Write([]byte(body))
			return
		}

		w.Header().Set("Cache-Control", "private, no-store")
		w.Header().Set("X-Cache", "BYPASS")
		w.Write([]byte(body))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
