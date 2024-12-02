package serve

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
)

func generateRandomBasicUsername() string {
	return gofakeit.Username()
}

func generateBasicPassword() string {
	length := 16
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func RunServer(port string) {
	username := generateRandomBasicUsername()
	fmt.Println("Username:", username)
	password := generateBasicPassword()
	fmt.Println("Password:", password)

	expectedUsernameHash := sha256.Sum256([]byte(username))
	expectedPasswordHash := sha256.Sum256([]byte(password))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
