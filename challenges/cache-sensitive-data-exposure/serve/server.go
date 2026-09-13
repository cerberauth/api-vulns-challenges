package serve

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type cache struct {
	mu    sync.Mutex
	store map[string]cachedResponse
}

type cachedResponse struct {
	body       string
	setCookie  string
	cacheable  bool
	csrfToken  string
	statusCode int
}

func newCache() *cache {
	return &cache{store: make(map[string]cachedResponse)}
}

func newToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func RunServer(port string, vulnerable bool) {
	c := newCache()

	// /dashboard issues a fresh session cookie and CSRF token per request
	// (as a login/dashboard page typically does), and returns a
	// personalized JSON payload. The vulnerable origin doesn't mark this
	// response as private, so a shared cache in front of it stores and
	// replays the first caller's session cookie, CSRF token and personal
	// data to everyone else who requests the same URL.
	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path

		c.mu.Lock()
		entry, hit := c.store[key]
		c.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if hit && entry.cacheable {
			w.Header().Set("Set-Cookie", entry.setCookie)
			w.Header().Set("X-CSRF-Token", entry.csrfToken)
			w.Header().Set("X-Cache", "HIT")
			w.Write([]byte(entry.body))
			return
		}

		sessionID := newToken()
		csrfToken := newToken()
		body := fmt.Sprintf(`{"user": "user-%s", "balance": 4200, "csrf_token": %q}`, sessionID[:6], csrfToken)
		setCookie := "session=" + sessionID + "; HttpOnly"

		entry = cachedResponse{
			body:      body,
			setCookie: setCookie,
			csrfToken: csrfToken,
			cacheable: vulnerable, // vulnerable: no Cache-Control: private, no-store is sent, so a shared cache stores it anyway
		}

		if !vulnerable {
			w.Header().Set("Cache-Control", "private, no-store")
		}

		if vulnerable {
			c.mu.Lock()
			c.store[key] = entry
			c.mu.Unlock()
		}

		w.Header().Set("Set-Cookie", setCookie)
		w.Header().Set("X-CSRF-Token", csrfToken)
		w.Header().Set("X-Cache", "MISS")
		w.Write([]byte(body))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
