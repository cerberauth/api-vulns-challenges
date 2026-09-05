package serve

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

func RunServer(port string, vulnerable bool) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(401)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			kid, ok := token.Header["kid"].(string)
			if !ok || kid == "" {
				return nil, fmt.Errorf("missing kid header")
			}

			keyPath := kid
			if !vulnerable {
				// fixed: the kid is restricted to a file name within the keys directory
				keyPath = filepath.Join("keys", filepath.Base(kid))
			}

			// VULNERABILITY: no path sanitization - allows path traversal
			keyBytes, err := os.ReadFile(keyPath)
			if err != nil {
				return nil, fmt.Errorf("key not found: %v", err)
			}

			return keyBytes, nil
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
