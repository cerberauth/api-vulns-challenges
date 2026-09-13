package serve

import (
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

// cacheKeyVulnerable mimics a CDN that normalizes path traversal segments
// BEFORE computing the cache key, while the origin below resolves the
// request by its own routing rules on the raw, un-normalized path - a
// discrepancy that lets a request for one path be stored under a
// completely different key.
//
// e.g. "/static/main.js/../../account" is normalized by the cache to
// "/account" (its cache key), but the origin's router matches the raw
// "/static/" prefix first and serves static-bundle content, so an
// attacker-chosen response ends up cached under the high-value "/account"
// key - a "cache what/where" primitive.
//
// Doyhenard, "Gotta cache 'em all" (Black Hat USA 2024).
func cacheKeyVulnerable(rawPath string) string {
	return normalizePath(rawPath)
}

// cacheKeyFixed keys on the raw, un-normalized path, matching how the
// origin actually routes it - no discrepancy, no collision.
func cacheKeyFixed(rawPath string) string {
	return rawPath
}

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

	// registered directly as the server's handler (not via http.ServeMux)
	// so the raw, un-normalized request path reaches the handler exactly as
	// sent, the same way a real reverse proxy would forward it upstream
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawPath := r.URL.Path

		var key string
		if vulnerable {
			key = cacheKeyVulnerable(rawPath)
		} else {
			key = cacheKeyFixed(rawPath)
		}

		c.mu.Lock()
		body, hit := c.store[key]
		c.mu.Unlock()

		w.Header().Set("Content-Type", "text/plain")
		if hit {
			w.Header().Set("X-Cache", "HIT")
			w.Write([]byte(body))
			return
		}

		// the origin routes on the raw path prefix, exactly like a real
		// reverse proxy / app router would, with no path normalization
		normalized := normalizePath(rawPath)
		switch {
		case strings.HasPrefix(rawPath, "/static/"):
			body = "// static bundle contents (attacker-influenced filename: " + rawPath + ")"
		case normalized == "/account":
			body = `{"balance": 1000, "csrf_token": "super-secret-csrf-token"}`
		default:
			w.WriteHeader(http.StatusNotFound)
			body = "not found"
		}

		c.mu.Lock()
		c.store[key] = body
		c.mu.Unlock()

		w.Header().Set("X-Cache", "MISS")
		w.Write([]byte(body))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
