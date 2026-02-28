package serve

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

func RunServer(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(401)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte("secret"), nil
		})

		if token != nil && token.Valid {
			w.WriteHeader(204)
		} else {
			fmt.Println(err)
			w.WriteHeader(401)
		}
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
