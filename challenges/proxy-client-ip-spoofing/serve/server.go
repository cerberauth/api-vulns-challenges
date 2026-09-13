package serve

import (
	"log"
	"net/http"
	"strings"
)

// trustedInternalIP is the only client IP allowed to reach the internal endpoint.
const trustedInternalIP = "127.0.0.1"

func clientIP(r *http.Request, vulnerable bool) string {
	if vulnerable {
		// vulnerable: the client-controlled X-Forwarded-For / X-Real-IP
		// headers are trusted as-is to make ACL and rate-limit decisions,
		// even though they can be freely spoofed by the caller
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			parts := strings.Split(xff, ",")
			return strings.TrimSpace(parts[0])
		}
		if xri := r.Header.Get("X-Real-IP"); xri != "" {
			return xri
		}
	}
	host := r.RemoteAddr
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}
	return host
}

func RunServer(port string, vulnerable bool) {
	http.HandleFunc("/internal/admin", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if clientIP(r, vulnerable) != trustedInternalIP {
			http.Error(w, `{"error": "forbidden"}`, http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "internal admin panel"}`))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
