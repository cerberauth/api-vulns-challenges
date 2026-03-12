package serve

import (
	"crypto"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

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

func RunServer(port string) {
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

		parts := strings.Split(tokenString, ".")
		sig, err := jwt.NewParser().DecodeSegment(parts[2])
		if err != nil {
			w.WriteHeader(401)
			return
		}

		err = jwt.GetSigningMethod(jwt.SigningMethodRS256.Alg()).Verify(strings.Join(parts[0:2], "."), sig, publicKey)
		if err == nil {
			w.WriteHeader(204)
		} else {
			fmt.Println(err)
			w.WriteHeader(401)
		}
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
