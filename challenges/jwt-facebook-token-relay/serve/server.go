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
	Issuer         = "https://www.facebook.com"
	VictimAudience = "1029384756192837"
)

func RunServer(port string, vulnerable bool) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	idpPublicKeyBytes, err := os.ReadFile(path.Join(cwd, "keys", "idp_public_key.pem"))
	if err != nil {
		log.Fatal(err)
	}

	idpPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(idpPublicKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(401)
			return
		}

		// The relying party checks the Facebook signature and issuer only.
		// It never validates the audience, so an ID token minted by Facebook
		// for a different, attacker-controlled app is accepted here too.
		parserOpts := []jwt.ParserOption{jwt.WithIssuer(Issuer)}
		if !vulnerable {
			parserOpts = append(parserOpts, jwt.WithAudience(VictimAudience))
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return idpPublicKey, nil
		}, parserOpts...)

		if err != nil || !token.Valid {
			fmt.Println(err)
			w.WriteHeader(401)
			return
		}

		w.WriteHeader(204)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
