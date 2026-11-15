package serve

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

const (
	DefaultKid    = "default"
	DefaultSecret = "Nq5Ty8Wd2Lm7Rx4Kv9Bc3Hp6Zs1Gf0Ja"
)

var kidPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func writeKeys() (string, error) {
	dir, err := os.MkdirTemp("", "keys")
	if err != nil {
		return "", err
	}
	return dir, os.WriteFile(filepath.Join(dir, DefaultKid+".key"), []byte(DefaultSecret), 0o600)
}

func vulnerableKey(dir, kid string) []byte {
	out, _ := exec.Command("sh", "-c", "cat "+dir+"/"+kid+".key").Output()
	return bytes.TrimSpace(out)
}

func fixedKey(dir, kid string) ([]byte, error) {
	if !kidPattern.MatchString(kid) {
		return nil, errors.New("invalid kid")
	}
	key, err := os.ReadFile(filepath.Join(dir, kid+".key"))
	if err != nil {
		return nil, err
	}
	if len(key) == 0 {
		return nil, errors.New("empty key")
	}
	return key, nil
}

func RunServer(port string, vulnerable bool) {
	dir, err := writeKeys()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			kid, ok := token.Header["kid"].(string)
			if !ok || kid == "" {
				return nil, errors.New("missing kid header")
			}
			if vulnerable {
				return vulnerableKey(dir, kid), nil
			}
			return fixedKey(dir, kid)
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
