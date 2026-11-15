package serve

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/subtle"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ClientID     = "challenge-client"
	ClientSecret = "challenge-client-secret"
	Username     = "alice"
	Password     = "wonderland"
	KeyID        = "challenge-key"
	TokenTTL     = 5 * time.Minute
)

const spec = `{
  "openapi": "3.0.3",
  "info": {"title": "JWT HMAC Confusion via JWKS", "version": "1.0.0"},
  "servers": [{"url": "http://localhost:%[1]s"}],
  "components": {
    "securitySchemes": {
      "oauth2": {
        "type": "oauth2",
        "flows": {
          "clientCredentials": {"tokenUrl": "http://localhost:%[1]s/oauth/token", "scopes": {"profile": "Read the profile"}},
          "password": {"tokenUrl": "http://localhost:%[1]s/oauth/token", "scopes": {"profile": "Read the profile"}}
        }
      }
    }
  },
  "security": [{"oauth2": ["profile"]}],
  "paths": {
    "/me": {"get": {"responses": {"200": {"description": "Profile"}, "401": {"description": "Unauthorized"}}}}
  }
}`

func equal(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func clientAuthenticated(r *http.Request) bool {
	id, secret, ok := r.BasicAuth()
	if !ok {
		id, secret = r.PostForm.Get("client_id"), r.PostForm.Get("client_secret")
	}
	return equal(id, ClientID) && equal(secret, ClientSecret)
}

func RunServer(port string, vulnerable bool) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		log.Fatal(err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
	issuer := "http://localhost:" + port

	jwks, err := json.Marshal(map[string]interface{}{"keys": []map[string]string{{
		"kty": "RSA",
		"use": "sig",
		"alg": "RS256",
		"kid": KeyID,
		"n":   base64.RawURLEncoding.EncodeToString(key.PublicKey.N.Bytes()),
		"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.PublicKey.E)).Bytes()),
	}}})
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, spec, port)
	})
	mux.HandleFunc("GET /.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jwks)
	})
	mux.HandleFunc("GET /.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"issuer":                                issuer,
			"token_endpoint":                        issuer + "/oauth/token",
			"jwks_uri":                              issuer + "/.well-known/jwks.json",
			"grant_types_supported":                 []string{"client_credentials", "password"},
			"id_token_signing_alg_values_supported": []string{"RS256"},
		})
	})
	mux.HandleFunc("POST /oauth/token", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, `{"error": "invalid_request"}`, http.StatusBadRequest)
			return
		}
		if !clientAuthenticated(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "invalid_client"}`))
			return
		}

		subject := ClientID
		switch r.PostForm.Get("grant_type") {
		case "client_credentials":
		case "password":
			if !equal(r.PostForm.Get("username"), Username) || !equal(r.PostForm.Get("password"), Password) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "invalid_grant"}`))
				return
			}
			subject = Username
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "unsupported_grant_type"}`))
			return
		}

		now := time.Now()
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"iss":   issuer,
			"sub":   subject,
			"scope": "profile",
			"iat":   now.Unix(),
			"exp":   now.Add(TokenTTL).Unix(),
		})
		token.Header["kid"] = KeyID
		accessToken, err := token.SignedString(key)
		if err != nil {
			http.Error(w, `{"error": "server_error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": accessToken,
			"token_type":   "Bearer",
			"expires_in":   int(TokenTTL.Seconds()),
			"scope":        "profile",
		})
	})
	mux.HandleFunc("GET /me", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok && vulnerable {
				return publicKeyPEM, nil
			}
			if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
				return &key.PublicKey, nil
			}
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}, jwt.WithExpirationRequired())
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sub, _ := claims.GetSubject()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"sub": sub})
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
