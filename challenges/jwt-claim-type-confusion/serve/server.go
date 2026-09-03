package serve

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

// HmacSecret is the API's dev-time signing key. It is intentionally known
// here so a valid token can be re-signed after its claims are tampered
// with - the bug this challenge demonstrates lives in how claim values are
// consumed, not in the signature check itself.
const HmacSecret = "s3cr3t-dev-key"

// debugRecoveryMiddleware is a leftover "helpful" dev-mode handler: rather
// than returning a generic error, it dumps the panic message and the full
// goroutine stack trace straight into the response body. It was meant to
// speed up local debugging and never got stripped out before this API
// shipped, so any request that crashes the handler leaks internals and
// returns a 500 with a body that looks nothing like a normal response.
func debugRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "panic: %v\n\n%s", rec, debug.Stack())
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RunServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(HmacSecret), nil
		})
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		// Both claims are trusted blindly and type-asserted straight to
		// the type the happy path expects, with no ok-check. A signature
		// check only proves who signed the token, not that its claim
		// values have the right shape - a claim of the wrong JSON type
		// (number, bool, object) or JSON null (which decodes to a bare
		// nil interface{} and fails the assertion exactly like a wrong
		// type does) crashes the handler instead of being rejected.
		name := claims["name"].(string)
		roles := claims["roles"].([]interface{})

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"name":%q,"roles":%d}`, name, len(roles))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(debugRecoveryMiddleware(mux))))
}
