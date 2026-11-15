package serve

import (
	"log"
	"net/http"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

const (
	HmacSecret = "claim-removal-dev-key"
	Issuer     = "https://idp.example.com"
	Audience   = "https://api.example.com"
)

func parserOptions(vulnerable bool) []jwt.ParserOption {
	opts := []jwt.ParserOption{jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()})}
	if vulnerable {
		return opts
	}
	return append(opts,
		jwt.WithIssuer(Issuer),
		jwt.WithAudience(Audience),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)
}

func validClaims(claims jwt.MapClaims, vulnerable bool) bool {
	if vulnerable {
		return true
	}
	sub, err := claims.GetSubject()
	if err != nil || sub == "" {
		return false
	}
	iat, err := claims.GetIssuedAt()
	return err == nil && iat != nil
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
		}, parserOptions(vulnerable)...)
		if err != nil || !token.Valid || !validClaims(claims, vulnerable) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
