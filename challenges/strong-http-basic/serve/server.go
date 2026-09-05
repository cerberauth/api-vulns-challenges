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
	"github.com/cerberauth/api-vulns-challenges/common"
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

func RunServer(port string, vulnerable bool) {
	var expectedUsername, expectedPassword string
	if vulnerable {
		// vulnerable: well-known, dictionary-guessable credentials
		expectedUsername = "admin"
		expectedPassword = "password"
	} else {
		expectedUsername = generateRandomBasicUsername()
		expectedPassword = generateBasicPassword()
	}
	fmt.Println("Username:", expectedUsername)
	fmt.Println("Password:", expectedPassword)

	expectedUsernameHash := sha256.Sum256([]byte(expectedUsername))
	expectedPasswordHash := sha256.Sum256([]byte(expectedPassword))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			var usernameMatch, passwordMatch bool
			if vulnerable {
				// vulnerable: plain, non-constant-time comparison, prone to timing attacks
				usernameMatch = username == expectedUsername
				passwordMatch = password == expectedPassword
			} else {
				usernameHash := sha256.Sum256([]byte(username))
				passwordHash := sha256.Sum256([]byte(password))

				usernameMatch = subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
				passwordMatch = subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1
			}

			if usernameMatch && passwordMatch {
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
