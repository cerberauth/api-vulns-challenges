package serve

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func RunServer() {
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
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})

		if token != nil && err == nil {
			w.WriteHeader(204)
		} else {
			fmt.Println(err)
			w.WriteHeader(401)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
