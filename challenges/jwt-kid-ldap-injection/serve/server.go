package serve

import (
	"errors"
	"log"
	"net/http"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

const (
	DefaultKid    = "default"
	DefaultSecret = "Hb6Zq1Xv8Mt3Pk5Wr2Lc9Nd4Fs7Gy0Ja"
)

var directory = []Entry{
	{"dn": "kid=rotation-pending,ou=keys,dc=example,dc=com", "objectclass": "jwtSigningKey", "kid": "rotation-pending", "uid": "svc-api", "secret": ""},
	{"dn": "kid=default,ou=keys,dc=example,dc=com", "objectclass": "jwtSigningKey", "kid": DefaultKid, "uid": "svc-api", "secret": DefaultSecret},
}

func lookupKey(kid string, vulnerable bool) ([]byte, error) {
	if !vulnerable {
		kid = EscapeFilterValue(kid)
	}
	match, err := ParseFilter("(&(objectClass=jwtSigningKey)(kid=" + kid + "))")
	if err != nil {
		return nil, err
	}
	for _, entry := range directory {
		if !match(entry) {
			continue
		}
		if !vulnerable && entry["secret"] == "" {
			return nil, errors.New("key disabled")
		}
		return []byte(entry["secret"]), nil
	}
	return nil, errors.New("key not found")
}

func RunServer(port string, vulnerable bool) {
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
			return lookupKey(kid, vulnerable)
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
