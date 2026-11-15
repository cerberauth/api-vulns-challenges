package main

import (
	"time"

	"github.com/cerberauth/api-vulns-challenges/challenges/jwt-scope-not-validated/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

func generateToken(_ bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   "2cb307ba-bb46-4194-854f-4774046d9c9b",
		"scope": "orders:read " + serve.RequiredScope,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString([]byte(serve.HmacSecret))
}

func main() {
	common.Execute(serve.RunServer, common.NewJwtCmd(generateToken))
}
