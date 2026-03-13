package main

import (
	"os"
	"path"
	"time"

	"github.com/cerberauth/api-vulns-challenges/challenges/jwt-kid-path-traversal/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

func generateToken() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	keyBytes, err := os.ReadFile(path.Join(cwd, "keys", "default.key"))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "2cb307ba-bb46-4194-854f-4774046d9c9b",
		"name": "John Doe",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	token.Header["kid"] = "keys/default.key"
	return token.SignedString(keyBytes)
}

func main() {
	common.Execute(serve.RunServer, common.NewJwtCmd(generateToken))
}
