package serve

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

var staticExtensions = []string{".css", ".js", ".png", ".jpg", ".ico", ".woff2"}

type cache struct {
	mu    sync.Mutex
	store map[string]string
}

func newCache() *cache {
	return &cache{store: make(map[string]string)}
}

func hasStaticExtension(path string) bool {
	for _, ext := range staticExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}

func RunServer(port string, vulnerable bool) {
	c := newCache()

	// /profile is a dynamic, authenticated endpoint. The shared cache in
	// front of it decides cacheability purely from the URL's file
	// extension, the classic web cache deception pattern: appending a
	// static-looking extension such as ".css" to a dynamic path fools the
	// cache into treating a per-user response as a shared static asset.
	//
	// The origin router below, in its vulnerable configuration, also
	// serves the same dynamic handler regardless of any trailing
	// extension, so the sensitive response is both produced AND cached
	// under a path anyone can guess and request without authentication.
	//
	// Omer Gil, original web cache deception research (2017).
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		serveProfile(w, r, c, vulnerable)
	})
	for _, ext := range staticExtensions {
		ext := ext
		http.HandleFunc("/profile"+ext, func(w http.ResponseWriter, r *http.Request) {
			serveProfile(w, r, c, vulnerable)
		})
	}

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func serveProfile(w http.ResponseWriter, r *http.Request, c *cache, vulnerable bool) {
	token := r.Header.Get("Authorization")

	if vulnerable && hasStaticExtension(r.URL.Path) {
		key := r.URL.Path
		c.mu.Lock()
		body, hit := c.store[key]
		c.mu.Unlock()
		if hit {
			w.Header().Set("X-Cache", "HIT")
			w.Write([]byte(body))
			return
		}
		body = fmt.Sprintf(`{"email": "%s@example.com", "session_token": "%s"}`, strings.TrimPrefix(token, "Bearer "), token)
		c.mu.Lock()
		c.store[key] = body
		c.mu.Unlock()
		w.Header().Set("X-Cache", "MISS")
		w.Write([]byte(body))
		return
	}

	if hasStaticExtension(r.URL.Path) {
		// fixed: the origin does not serve dynamic, authenticated content
		// under any path carrying a static-looking extension
		w.Header().Set("Cache-Control", "no-store")
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Cache", "BYPASS")
	body := fmt.Sprintf(`{"email": "%s@example.com", "session_token": "%s"}`, strings.TrimPrefix(token, "Bearer "), token)
	w.Write([]byte(body))
}
