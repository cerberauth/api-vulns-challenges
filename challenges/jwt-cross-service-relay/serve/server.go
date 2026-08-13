package serve

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

const (
	Issuer           = "https://idp.example.com"
	ServiceBAudience = "service-b-client-id"
)

func RunServer(port string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	idpPublicKeyBytes, err := os.ReadFile(path.Join(cwd, "keys", "idp_public_key.pem"))
	if err != nil {
		log.Fatal(err)
	}

	idpPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(idpPublicKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(401)
			return
		}

		// Service B only checks the signature and issuer of the shared IdP.
		// It never validates the audience, so any token minted by the IdP
		// for another service (e.g. Service A) is accepted here too.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return idpPublicKey, nil
		}, jwt.WithIssuer(Issuer))

		if err != nil || !token.Valid {
			fmt.Println(err)
			w.WriteHeader(401)
			return
		}

		w.WriteHeader(204)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
