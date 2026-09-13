package serve

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
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

// decodeOnce percent-decodes a path a single time, mimicking a cache that
// stops after one decoding pass. A nested/double-encoded segment (e.g.
// "%252F", which is "%2F" with its own "%" escaped) survives a single pass
// as the still-opaque literal text "%2F", hiding a real "/" path separator
// from the cache's view of the path.
func decodeOnce(path string) string {
	decoded, err := url.PathUnescape(path)
	if err != nil {
		return path
	}
	return decoded
}

// decodeFully keeps percent-decoding until no escapes remain, then collapses
// "." / ".." segments - the same normalization an origin router applies
// before matching a route.
func decodeFully(path string) string {
	for {
		decoded, err := url.PathUnescape(path)
		if err != nil || decoded == path {
			break
		}
		path = decoded
	}
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

	// registered directly (bypassing http.ServeMux) so encoded path
	// segments reach the handler unmodified
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawPath := r.RequestURI
		if idx := strings.Index(rawPath, "?"); idx != -1 {
			rawPath = rawPath[:idx]
		}
		token := r.Header.Get("Authorization")

		// the origin always fully decodes and normalizes before routing
		originPath := decodeFully(rawPath)

		cacheableStatic := false
		if vulnerable {
			// the cache decodes only once, so a double-encoded "/" is
			// still hidden inside what looks like a single opaque
			// "/static/..." segment
			cacheableStatic = strings.HasPrefix(decodeOnce(rawPath), "/static/")
		} else {
			// fixed: the cache decodes and normalizes exactly like the
			// origin does, so it sees the same effective path
			cacheableStatic = strings.HasPrefix(originPath, "/static/")
		}

		var body string
		switch {
		case originPath == "/account":
			body = fmt.Sprintf(`{"email": "%s@example.com", "session_token": "%s"}`, strings.TrimPrefix(token, "Bearer "), token)
		default:
			body = "// static asset contents"
		}

		if cacheableStatic && originPath == "/account" {
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

		if originPath == "/account" {
			w.Header().Set("Cache-Control", "private, no-store")
			w.Header().Set("X-Cache", "BYPASS")
		}
		w.Write([]byte(body))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
