package serve

import (
	"crypto"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

func readPublicKey() (crypto.PublicKey, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := os.ReadFile(path.Join(cwd, "keys", "public_key.pem"))
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
}

func RunServer(port string, vulnerable bool) {
	publicKey, err := readPublicKey()
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

		var valid bool
		if vulnerable {
			// vulnerable: the signature is never checked
			token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
			valid = token != nil && err == nil
		} else {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return publicKey, nil
			})
			valid = err == nil && token.Valid
		}

		if valid {
			w.WriteHeader(204)
		} else {
			w.WriteHeader(401)
		}
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
