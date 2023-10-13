package serve

import (
	"crypto"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func readPublicKey() (crypto.PublicKey, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := os.ReadFile(cwd + string(os.PathSeparator) + "keys/public_key.pem")
	if err != nil {
		return nil, err
	}

	return jwt.ParseEdPublicKeyFromPEM(publicKeyBytes)
}

func RunServer() {
	publicKey, err := readPublicKey()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader == "" {
			w.WriteHeader(401)
			return
		}

		headerParts := strings.Split(authorizationHeader, "Bearer")
		if len(headerParts) != 2 {
			w.WriteHeader(401)
			return
		}

		tokenString := strings.TrimSpace(headerParts[1])
		parts := strings.Split(tokenString, ".")
		sig, err := jwt.NewParser().DecodeSegment(parts[2])
		if err != nil {
			w.WriteHeader(401)
			return
		}

		err = jwt.GetSigningMethod("EdDSA").Verify(strings.Join(parts[0:2], "."), sig, publicKey)
		if err == nil {
			w.WriteHeader(204)
		} else {
			fmt.Println(err)
			w.WriteHeader(401)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
