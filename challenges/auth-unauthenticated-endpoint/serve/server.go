package serve

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

const Secret = "pR8vQw2LxN5tYz7KcB4mJf6HsD9gUa3E"

const spec = `{
  "openapi": "3.0.3",
  "info": {"title": "Unauthenticated Endpoint", "version": "1.0.0"},
  "servers": [{"url": "http://localhost:%s"}],
  "components": {
    "securitySchemes": {"bearerAuth": {"type": "http", "scheme": "bearer", "bearerFormat": "JWT"}}
  },
  "paths": {
    "/account": {
      "get": {
        "security": [{"bearerAuth": []}],
        "responses": {"200": {"description": "Account"}, "302": {"description": "Login redirect"}}
      }
    },
    "/status": {
      "get": {
        "security": [{"bearerAuth": []}, {}],
        "responses": {"200": {"description": "Status"}}
      }
    }
  }
}`

const loginPage = `<!DOCTYPE html><html><head><title>Sign in</title></head><body><form method="post" action="/login"><input name="username"><input name="password" type="password"><button>Sign in</button></form></body></html>`

const appPage = `<!DOCTYPE html><html><head><title>App</title></head><body><div id="root"></div></body></html>`

func authenticated(r *http.Request) bool {
	tokenString, ok := common.ExtractBearerToken(r)
	if !ok {
		return false
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	return err == nil && token.Valid
}

func writeHTML(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func RunServer(port string, vulnerable bool) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, spec, port)
	})
	mux.HandleFunc("GET /account", func(w http.ResponseWriter, r *http.Request) {
		if !vulnerable && !authenticated(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id": "2cb307ba-bb46-4194-854f-4774046d9c9b", "email": "john.doe@example.com", "plan": "premium"}`))
	})
	mux.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
	})
	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		writeHTML(w, loginPage)
	})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		writeHTML(w, appPage)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
