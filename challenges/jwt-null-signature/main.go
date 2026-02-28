package main

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/cerberauth/api-vulns-challenges/challenges/jwt-null-signature/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

func generateToken() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	privateKeyBytes, err := os.ReadFile(path.Join(cwd, "keys", "private_key.pem"))
	if err != nil {
		return "", err
	}

	key, err := jwt.ParseEdPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", err
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"sub":  "2cb307ba-bb46-4194-854f-4774046d9c9b",
		"name": "John Doe",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour).Unix(),
	}).SignedString(key)
	if err != nil {
		return "", err
	}

	parts := strings.Split(tokenString, ".")
	return strings.Join([]string{parts[0], parts[1], ""}, "."), nil
}

func main() {
	common.Execute(serve.RunServer, common.NewJwtCmd(generateToken))
}
