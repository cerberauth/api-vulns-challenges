package common

import (
	"net/http"
	"strings"
)

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
