package main

import (
	"time"

	"github.com/cerberauth/api-vulns-challenges/challenges/jwt-missing-expiry/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

func generateToken(vulnerable bool) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": "2cb307ba-bb46-4194-854f-4774046d9c9b",
		"iat": now.Unix(),
		"exp": now.Add(serve.TokenLifetime).Unix(),
	}
	if vulnerable {
		delete(claims, "exp")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(serve.HmacSecret))
}

func main() {
	common.Execute(serve.RunServer, common.NewJwtCmd(generateToken))
}
