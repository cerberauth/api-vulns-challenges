package serve

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	maxAge               = 2 * time.Second
	staleWhileRevalidate = 5 * time.Second
	staleIfError         = 5 * time.Second
)

type cacheEntry struct {
	body      string
	storedAt  time.Time
	revalDone bool
}

type cache struct {
	mu    sync.Mutex
	entry *cacheEntry
}

// failNow controls whether the origin is currently simulating an outage,
// toggled via /toggle-outage so a client can exercise stale-if-error.
type originState struct {
	mu   sync.Mutex
	down bool
	gen  int
}

func RunServer(port string, vulnerable bool) {
	c := &cache{}
	origin := &originState{}

	fetchFromOrigin := func() (string, int, bool) {
		origin.mu.Lock()
		down := origin.down
		origin.gen++
		gen := origin.gen
		origin.mu.Unlock()
		if down {
			return "", http.StatusInternalServerError, false
		}
		return time.Now().Format(time.RFC3339Nano) + " gen=" + strconv.Itoa(gen), http.StatusOK, true
	}

	http.HandleFunc("/toggle-outage", func(w http.ResponseWriter, r *http.Request) {
		origin.mu.Lock()
		origin.down = !origin.down
		down := origin.down
		origin.mu.Unlock()
		if down {
			w.Write([]byte("origin is now DOWN"))
		} else {
			w.Write([]byte("origin is now UP"))
		}
	})

	// / declares stale-while-revalidate and stale-if-error, and reports
	// Cache-Control on every response so a scanner can confirm the
	// directives are actually honored by the fronting cache (this handler
	// itself plays the role of the shared cache + origin combined).
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Cache-Control", "max-age=2, stale-while-revalidate=5, stale-if-error=5")

		c.mu.Lock()
		entry := c.entry
		c.mu.Unlock()

		now := time.Now()
		if entry != nil {
			age := now.Sub(entry.storedAt)
			fresh := age <= maxAge

			if fresh {
				w.Header().Set("X-Cache", "HIT")
				w.Write([]byte(entry.body))
				return
			}

			if vulnerable {
				// vulnerable: the cache never serves stale content and
				// never honors stale-if-error - it always goes back to
				// the origin once max-age has elapsed, defeating both
				// directives even though it advertises them
				body, status, ok := fetchFromOrigin()
				if !ok {
					http.Error(w, "origin unavailable", status)
					return
				}
				c.mu.Lock()
				c.entry = &cacheEntry{body: body, storedAt: now}
				c.mu.Unlock()
				w.Header().Set("X-Cache", "MISS")
				w.Write([]byte(body))
				return
			}

			// fixed: within the stale-while-revalidate window, serve the
			// stale entry immediately and refresh it in the background
			if age <= maxAge+staleWhileRevalidate {
				w.Header().Set("X-Cache", "STALE")
				w.Write([]byte(entry.body))
				go func() {
					body, _, ok := fetchFromOrigin()
					if ok {
						c.mu.Lock()
						c.entry = &cacheEntry{body: body, storedAt: time.Now()}
						c.mu.Unlock()
					}
				}()
				return
			}

			// fixed: within the stale-if-error window, serve the stale
			// entry if the origin errors out instead of surfacing the error
			if age <= maxAge+staleIfError {
				body, status, ok := fetchFromOrigin()
				if !ok {
					_ = status
					w.Header().Set("X-Cache", "STALE")
					w.Write([]byte(entry.body))
					return
				}
				c.mu.Lock()
				c.entry = &cacheEntry{body: body, storedAt: now}
				c.mu.Unlock()
				w.Header().Set("X-Cache", "MISS")
				w.Write([]byte(body))
				return
			}
		}

		body, status, ok := fetchFromOrigin()
		if !ok {
			http.Error(w, "origin unavailable", status)
			return
		}
		c.mu.Lock()
		c.entry = &cacheEntry{body: body, storedAt: now}
		c.mu.Unlock()
		w.Header().Set("X-Cache", "MISS")
		w.Write([]byte(body))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
