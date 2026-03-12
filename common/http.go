package common

import (
	"net/http"
	"strings"
)

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; sandbox")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}

func ExtractBearerToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("authorization")
	if authHeader == "" {
		return "", false
	}
	parts := strings.Split(authHeader, "Bearer")
	if len(parts) != 2 {
		return "", false
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}
	return token, true
}
