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

func RunServer(port string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	publicKeyBytes, err := os.ReadFile(path.Join(cwd, "keys", "public_key.pem"))
	if err != nil {
		log.Fatal(err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/public-key", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-pem-file")
		w.Write(publicKeyBytes)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(401)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// VULNERABILITY: accepts HMAC tokens using RSA public key as secret
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
				return publicKeyBytes, nil
			}

			if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
				return publicKey, nil
			}

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		})

		if token != nil && token.Valid {
			w.WriteHeader(204)
		} else {
			fmt.Println(err)
			w.WriteHeader(401)
		}
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
