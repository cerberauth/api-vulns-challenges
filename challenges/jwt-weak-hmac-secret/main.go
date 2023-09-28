package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader == "" {
			w.WriteHeader(401)
			return
		}

		parts := strings.Split(authorizationHeader, "Bearer")
		if len(parts) != 2 {
			w.WriteHeader(401)
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
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

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
