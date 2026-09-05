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
// with - the bug this challenge demonstrates lives in how a claim value is
// consumed, not in the signature check itself.
const HmacSecret = "s3cr3t-dev-key"

// discountTiers maps a coupon code's length to a discount tier. Nobody
// designing this ever expected a coupon code longer than a few characters,
// so the slice was sized generously "just in case" and never bounds-checked.
var discountTiers = []string{
	"0%", "0%", "0%", "0%", "0%", "5%", "5%", "5%", "10%", "10%",
	"10%", "15%", "15%", "15%", "20%", "20%", "20%", "25%", "25%", "25%",
	"30%", "30%", "30%", "35%", "35%", "35%", "40%", "40%", "40%", "45%",
	"45%", "45%",
}

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
	mux.HandleFunc("/checkout", func(w http.ResponseWriter, r *http.Request) {
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

		var coupon string
		if vulnerable {
			coupon = claims["coupon"].(string)
		} else {
			coupon, ok = claims["coupon"].(string)
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		var discount string
		if vulnerable {
			// The coupon's raw length indexes straight into a fixed-size
			// lookup table with no upper bound. Every coupon the developers
			// ever tested with was a handful of characters, so this always
			// worked - until a claim value far longer than expected (which
			// is exactly what a length/size-boundary fuzz mutation sends)
			// walks past the end of the slice.
			discount = discountTiers[len(coupon)]
		} else {
			// fixed: the index is clamped to the table's bounds
			idx := len(coupon)
			if idx >= len(discountTiers) {
				idx = len(discountTiers) - 1
			}
			discount = discountTiers[idx]
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"coupon":%q,"discount":%q}`, coupon, discount)
	})

	log.Println("Server started at port", port)
	recovery := debugRecoveryMiddleware
	if !vulnerable {
		recovery = recoveryMiddleware
	}
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(recovery(mux))))
}
