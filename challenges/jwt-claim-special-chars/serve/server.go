package serve

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

// HmacSecret is the API's dev-time signing key. It is intentionally known
// here so a valid token can be re-signed after its claims are tampered
// with - the bug this challenge demonstrates lives in how a claim value is
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

// recoveryMiddleware is the fixed counterpart of debugRecoveryMiddleware: it
// still recovers from panics so the server keeps running, but never leaks
// the panic message or stack trace to the client.
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Println("recovered from panic:", rec)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RunServer(port string, vulnerable bool) {
	mux := http.NewServeMux()
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
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

		var rawFilter string
		if vulnerable {
			rawFilter = claims["filter"].(string)
		} else {
			rawFilter, ok = claims["filter"].(string)
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		var filter *regexp.Regexp
		if vulnerable {
			// The "filter" claim lets a client scope which of their own
			// records get returned, so it's compiled straight into a regexp
			// with MustCompile instead of Compile+error-check - nobody
			// expected a claim to contain anything but a simple pattern.
			// Special/meta characters that don't form valid regexp syntax
			// (unbalanced brackets, braces, parens, a trailing backslash, ...)
			// panic instead of being rejected as a bad filter.
			filter = regexp.MustCompile(rawFilter)
		} else {
			// fixed: invalid patterns are rejected instead of panicking
			var err error
			filter, err = regexp.Compile(rawFilter)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		records := []string{"invoice-1001", "invoice-1002", "invoice-1003"}
		matches := make([]string, 0, len(records))
		for _, record := range records {
			if filter.MatchString(record) {
				matches = append(matches, record)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"filter":%q,"matches":%d}`, rawFilter, len(matches))
	})

	log.Println("Server started at port", port)
	recovery := debugRecoveryMiddleware
	if !vulnerable {
		recovery = recoveryMiddleware
	}
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(recovery(mux))))
}
