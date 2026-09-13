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

// staticDirPrefixes are directory prefixes that the shared cache always
// caches, regardless of content-type, because the CDN rule is written as
// "cache everything under /static/**".
var staticDirPrefixes = []string{"/static/", "/assets/"}

func normalizePath(path string) string {
	segments := strings.Split(path, "/")
	var stack []string
	for _, s := range segments {
		switch s {
		case "", ".":
			continue
		case "..":
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		default:
			stack = append(stack, s)
		}
	}
	return "/" + strings.Join(stack, "/")
}

func RunServer(port string, vulnerable bool) {
	c := newCache()

	// registered directly (bypassing http.ServeMux) so path traversal
	// segments in the raw path reach the handler unmodified
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawPath := r.URL.Path
		token := r.Header.Get("Authorization")

		cacheableByRule := false
		if vulnerable {
			for _, prefix := range staticDirPrefixes {
				if strings.HasPrefix(rawPath, prefix) {
					// the CDN caches anything whose RAW path starts under
					// a static directory, without normalizing traversal
					// segments first
					cacheableByRule = true
					break
				}
			}
		}

		// the origin normalizes the path and routes on the *normalized*
		// result, so "/static/x/../../account" is served as the dynamic
		// /account handler even though the cache saw a "/static/" prefix
		normalized := normalizePath(rawPath)

		var body string
		switch {
		case normalized == "/account":
			body = fmt.Sprintf(`{"email": "%s@example.com", "session_token": "%s"}`, strings.TrimPrefix(token, "Bearer "), token)
		default:
			body = "// static asset contents"
		}

		if cacheableByRule && normalized == "/account" {
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

		if normalized == "/account" {
			w.Header().Set("Cache-Control", "private, no-store")
			w.Header().Set("X-Cache", "BYPASS")
		}
		w.Write([]byte(body))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
