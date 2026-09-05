package main

import (
	"time"

	"github.com/cerberauth/api-vulns-challenges/challenges/jwt-alg-none-bypass/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

func generateToken(vulnerable bool) (string, error) {
	claims := jwt.MapClaims{
		"sub":  "2cb307ba-bb46-4194-854f-4774046d9c9b",
		"name": "John Doe",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour).Unix(),
	}
	if vulnerable {
		token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		return token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(serve.Secret))
}

func main() {
	common.Execute(serve.RunServer, common.NewJwtCmd(generateToken))
}
