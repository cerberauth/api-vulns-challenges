package serve

import (
	"log"
	"net/http"
	"time"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

const (
	HmacSecret    = "Xk4mW9qT2vB7nL5rH8cF3jD6yP1sG0zE"
	TokenLifetime = 15 * time.Minute
	MaxLifetime   = time.Hour
)

func lifetimeAllowed(claims jwt.MapClaims) bool {
	exp, err := claims.GetExpirationTime()
	if err != nil || exp == nil {
		return false
	}
	iat, err := claims.GetIssuedAt()
	if err != nil || iat == nil {
		return false
	}
	return exp.Sub(iat.Time) <= MaxLifetime
}

func RunServer(port string, vulnerable bool) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(HmacSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
		if err != nil || !token.Valid || (!vulnerable && !lifetimeAllowed(claims)) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
