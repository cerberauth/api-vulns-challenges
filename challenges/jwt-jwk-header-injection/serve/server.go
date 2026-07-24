package serve

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

const LegitKid = "legit-key-1"

// jwkFromHeader turns the raw "jwk" header claim (decoded from JSON as
// map[string]interface{}) into an RSA public key, with no check of any kind
// against the server's own trusted key or a known keyset.
func jwkFromHeader(raw interface{}) (*rsa.PublicKey, error) {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid jwk header: not an object")
	}

	nStr, _ := m["n"].(string)
	eStr, _ := m["e"].(string)
	if nStr == "" || eStr == "" {
		return nil, fmt.Errorf("invalid jwk header: missing n or e")
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, fmt.Errorf("invalid jwk header: bad n: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, fmt.Errorf("invalid jwk header: bad e: %w", err)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: int(new(big.Int).SetBytes(eBytes).Int64()),
	}, nil
}

func RunServer(port string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	legitPublicKeyBytes, err := os.ReadFile(path.Join(cwd, "keys", "legit_public_key.pem"))
	if err != nil {
		log.Fatal(err)
	}

	legitPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(legitPublicKeyBytes)
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

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// VULNERABLE (CVE-2018-0114): a self-declared "jwk" header is
			// trusted as the verification key instead of being checked
			// against the server's own pinned key. Any token can carry its
			// own key and sign itself with it.
			if rawJWK, present := token.Header["jwk"]; present {
				return jwkFromHeader(rawJWK)
			}

			return legitPublicKey, nil
		})

		if err == nil && token.Valid {
			w.WriteHeader(204)
		} else {
			fmt.Println(err)
			w.WriteHeader(401)
		}
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
