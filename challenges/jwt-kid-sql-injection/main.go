package main

import (
	"time"

	"github.com/cerberauth/api-vulns-challenges/challenges/jwt-kid-sql-injection/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

const defaultKid = "default"
const defaultSecret = "supersecretkey_stored_in_database"

func generateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "2cb307ba-bb46-4194-854f-4774046d9c9b",
		"name": "John Doe",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	token.Header["kid"] = defaultKid
	return token.SignedString([]byte(defaultSecret))
}

func main() {
	common.Execute(serve.RunServer, common.NewJwtCmd(generateToken))
}
