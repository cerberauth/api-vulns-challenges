package main

import (
	"time"

	"github.com/cerberauth/api-vulns-challenges/challenges/jwt-timestamp-tampering/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

func generateToken(vulnerable bool) (string, error) {
	// Same freshly-issued token in both modes: iat is now, exp is one hour
	// out. The two server modes differ in whether they still enforce exp
	// (and nbf, if present) once the signature checks out - see
	// serve/server.go.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "2cb307ba-bb46-4194-854f-4774046d9c9b",
		"name": "John Doe",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString([]byte(serve.HmacSecret))
}

func main() {
	common.Execute(serve.RunServer, common.NewJwtCmd(generateToken))
}
