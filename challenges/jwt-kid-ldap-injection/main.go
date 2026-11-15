package main

import (
	"time"

	"github.com/cerberauth/api-vulns-challenges/challenges/jwt-kid-ldap-injection/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

func generateToken(_ bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "2cb307ba-bb46-4194-854f-4774046d9c9b",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	token.Header["kid"] = serve.DefaultKid
	return token.SignedString([]byte(serve.DefaultSecret))
}

func main() {
	common.Execute(serve.RunServer, common.NewJwtCmd(generateToken))
}
