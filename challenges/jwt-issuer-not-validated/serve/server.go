package serve

import (
	"log"
	"net/http"
	"strings"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

const (
	IdPSecret   = "Vt7Kx2Qm9Rb4Wn6Lp1Zc8Hs3Jd5Fg0Ya"
	IdPBaseURL  = "https://idp.example.com/"
	Issuer      = IdPBaseURL + "tenant-a"
	OtherIssuer = IdPBaseURL + "tenant-b"
)

func RunServer(port string, vulnerable bool) {
	opts := []jwt.ParserOption{jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired()}
	if !vulnerable {
		opts = append(opts, jwt.WithIssuer(Issuer))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(IdPSecret), nil
		}, opts...)
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if iss, _ := claims.GetIssuer(); !strings.HasPrefix(iss, IdPBaseURL) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
