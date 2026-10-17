package serve

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

// HmacSecret is the API's dev-time signing key. It is intentionally known
// here so it can be recovered by a secret-cracking tool and used to re-sign
// a tampered token - the bug this challenge demonstrates lives in whether
// the exp/nbf claims are actually enforced, not in the signature check
// itself.
const HmacSecret = "s3cr3t-dev-key"

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(HmacSecret), nil
}

func RunServer(port string, vulnerable bool) {
	mux := http.NewServeMux()
	mux.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var token *jwt.Token
		var err error
		if vulnerable {
			// vulnerable: claims validation (exp, nbf, iat, ...) is
			// explicitly turned off - only the signature is checked. A
			// token that is expired per its own exp claim, or one
			// re-signed with exp pushed arbitrarily far into the future,
			// is accepted exactly like a fresh one.
			token, err = jwt.NewParser(jwt.WithoutClaimsValidation()).Parse(tokenString, keyFunc)
		} else {
			// fixed: the default parser validates exp/nbf/iat as usual.
			token, err = jwt.Parse(tokenString, keyFunc)
		}
		if err != nil || token == nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
