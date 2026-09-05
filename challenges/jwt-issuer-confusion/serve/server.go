package serve

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

const (
	IssuerA = "https://issuer-a.example.com"
	IssuerB = "https://issuer-b.example.com"
)

func RunServer(port string, vulnerable bool) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	issuerAPublicKeyBytes, err := os.ReadFile(path.Join(cwd, "keys", "issuer-a_public_key.pem"))
	if err != nil {
		log.Fatal(err)
	}

	issuerAPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(issuerAPublicKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	issuerBPublicKeyBytes, err := os.ReadFile(path.Join(cwd, "keys", "issuer-b_public_key.pem"))
	if err != nil {
		log.Fatal(err)
	}

	issuerBPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(issuerBPublicKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	parser := jwt.NewParser()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(401)
			return
		}

		unverifiedClaims := jwt.MapClaims{}
		if _, _, err := parser.ParseUnverified(tokenString, unverifiedClaims); err != nil {
			w.WriteHeader(401)
			return
		}

		issuer, err := unverifiedClaims.GetIssuer()
		if err != nil {
			w.WriteHeader(401)
			return
		}

		var token *jwt.Token
		switch issuer {
		case IssuerA:
			token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return issuerAPublicKey, nil
			})
		case IssuerB:
			if vulnerable {
				// VULNERABILITY: issuer-b is a legacy/trusted issuer whose tokens are
				// parsed without any signature verification, so claims (including "iss")
				// can be forged freely as long as the "iss" value matches this branch.
				token, _, err = parser.ParseUnverified(tokenString, jwt.MapClaims{})
				if err == nil {
					token.Valid = true
				}
				break
			}

			token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return issuerBPublicKey, nil
			})
		default:
			err = fmt.Errorf("unknown issuer: %v", issuer)
		}

		if token != nil && token.Valid {
			w.WriteHeader(204)
		} else {
			fmt.Println(err)
			w.WriteHeader(401)
		}
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
