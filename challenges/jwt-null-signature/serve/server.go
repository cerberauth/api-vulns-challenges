package serve

import (
	"crypto"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

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

	return jwt.ParseEdPublicKeyFromPEM(publicKeyBytes)
}

func RunServer(port string) {
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
		if len(tokenString) == 0 {
			w.WriteHeader(401)
			return
		}

		parts := strings.Split(tokenString, ".")
		if len(parts[2]) == 0 {
			w.WriteHeader(204)
			return
		}

		sig, err := jwt.NewParser().DecodeSegment(parts[2])
		if err != nil {
			w.WriteHeader(401)
			return
		}

		err = jwt.GetSigningMethod(jwt.SigningMethodEdDSA.Alg()).Verify(strings.Join(parts[0:2], "."), sig, publicKey)
		if err == nil {
			w.WriteHeader(204)
		} else {
			fmt.Println(err)
			w.WriteHeader(401)
		}
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
