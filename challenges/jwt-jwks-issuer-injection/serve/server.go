package serve

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

const (
	LegitIssuer = "http://localhost:8090"
	LegitKid    = "legit-key-1"
	idpPort     = "8090"
)

type jwk struct {
	Kty string `json:"kty"`
	Use string `json:"use,omitempty"`
	Kid string `json:"kid"`
	Alg string `json:"alg,omitempty"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwks struct {
	Keys []jwk `json:"keys"`
}

func rsaPublicKeyToJWK(key *rsa.PublicKey, kid string) jwk {
	return jwk{
		Kty: "RSA",
		Use: "sig",
		Kid: kid,
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
	}
}

func jwkToRSAPublicKey(k jwk) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
	if err != nil {
		return nil, err
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
	if err != nil {
		return nil, err
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: int(new(big.Int).SetBytes(eBytes).Int64()),
	}, nil
}

// fetchJWKS builds the JWKS endpoint straight from the token's own issuer
// claim and fetches it, with no allowlist of trusted issuers.
func fetchJWKS(issuer string) (*jwks, error) {
	resp, err := http.Get(strings.TrimRight(issuer, "/") + "/.well-known/jwks.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jwks endpoint returned status %d", resp.StatusCode)
	}

	var set jwks
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		return nil, err
	}
	return &set, nil
}

// runIdentityProvider serves the legitimate issuer's JWKS document, so the
// challenge is runnable standalone without any external identity provider.
func runIdentityProvider(publicKey *rsa.PublicKey) {
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwks{Keys: []jwk{rsaPublicKeyToJWK(publicKey, LegitKid)}})
	})

	log.Println("Legit identity provider JWKS endpoint started at port", idpPort)
	log.Fatal(http.ListenAndServe(":"+idpPort, mux))
}

func RunServer(port string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	legitPublicKeyBytes, err := os.ReadFile(path.Join(cwd, "keys", "legit_public_key.pem"))
	if err != nil {
		log.Fatal(err)
	}

	legitPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(legitPublicKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	go runIdentityProvider(legitPublicKey)

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

		set, err := fetchJWKS(issuer)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(401)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			kid, _ := token.Header["kid"].(string)
			for _, key := range set.Keys {
				if key.Kid == kid {
					return jwkToRSAPublicKey(key)
				}
			}
			return nil, fmt.Errorf("no matching key found for kid: %v", kid)
		})

		if err == nil && token.Valid {
			w.WriteHeader(204)
		} else {
			fmt.Println(err)
			w.WriteHeader(401)
		}
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
