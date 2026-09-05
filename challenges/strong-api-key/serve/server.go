package serve

import (
	crand "crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"

	"github.com/cerberauth/api-vulns-challenges/common"
)

func generateStrongAPIKey() (string, error) {
	b := make([]byte, 32)
	_, err := crand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// generateWeakAPIKey mimics a short, predictable key derived from a
// low-entropy, non-cryptographic random source - practical to brute-force.
func generateWeakAPIKey() string {
	const digits = "0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = digits[rand.IntN(len(digits))]
	}
	return string(b)
}

func RunServer(port string, vulnerable bool) {
	var apiKey string
	if vulnerable {
		apiKey = generateWeakAPIKey()
	} else {
		var err error
		apiKey, err = generateStrongAPIKey()
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("API Key:", apiKey)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		provided := r.Header.Get("X-API-Key")

		var valid bool
		if vulnerable {
			// vulnerable: non-constant-time comparison, in addition to the
			// key itself being short and guessable
			valid = provided == apiKey
		} else {
			valid = subtle.ConstantTimeCompare([]byte(provided), []byte(apiKey)) == 1
		}

		if !valid {
			w.WriteHeader(401)
			return
		}

		w.WriteHeader(204)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
