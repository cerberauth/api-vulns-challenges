package serve

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	limitCount  = 5
	limitWindow = 10 * time.Second
)

type limiter struct {
	mu   sync.Mutex
	hits map[string][]time.Time
}

func newLimiter() *limiter {
	return &limiter{hits: make(map[string][]time.Time)}
}

func (l *limiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-limitWindow)

	kept := l.hits[key][:0]
	for _, t := range l.hits[key] {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	l.hits[key] = kept

	if len(l.hits[key]) >= limitCount {
		return false
	}
	l.hits[key] = append(l.hits[key], now)
	return true
}

// rateLimitKey identifies the caller for throttling purposes.
func rateLimitKey(r *http.Request, vulnerable bool) string {
	if vulnerable {
		// vulnerable: the rate limit key is taken from the client-controlled
		// X-Forwarded-For header, so an attacker can bypass the limit simply
		// by sending a different value on every request
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			return xff
		}
	}
	host := r.RemoteAddr
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}
	return host
}

func RunServer(port string, vulnerable bool) {
	l := newLimiter()

	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if !l.allow(rateLimitKey(r, vulnerable)) {
			http.Error(w, `{"error": "too many requests"}`, http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "login attempt processed"}`))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
