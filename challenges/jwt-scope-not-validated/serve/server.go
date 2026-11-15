package serve

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

const (
	HmacSecret    = "scope-dev-key"
	RequiredScope = "orders:write"
)

func hasScope(claims jwt.MapClaims, scope string) bool {
	value, ok := claims["scope"].(string)
	return ok && slices.Contains(strings.Fields(value), scope)
}

func hasIdentity(claims jwt.MapClaims) bool {
	sub, err := claims.GetSubject()
	if err != nil || sub == "" {
		return false
	}
	iat, err := claims.GetIssuedAt()
	return err == nil && iat != nil
}

func RunServer(port string, vulnerable bool) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(HmacSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !vulnerable && !hasIdentity(claims) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !vulnerable && !hasScope(claims, RequiredScope) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		var order map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": "ord_1", "status": "created"}`))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
