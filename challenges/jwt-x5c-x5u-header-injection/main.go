package main

import (
	"encoding/base64"
	"encoding/pem"
	"os"
	"path"
	"time"

	"github.com/cerberauth/api-vulns-challenges/challenges/jwt-x5c-x5u-header-injection/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

func generateToken() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	privateKeyBytes, err := os.ReadFile(path.Join(cwd, "keys", "legit_private_key.pem"))
	if err != nil {
		return "", err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", err
	}

	certBytes, err := os.ReadFile(path.Join(cwd, "keys", "legit_cert.pem"))
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(certBytes)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":  "2cb307ba-bb46-4194-854f-4774046d9c9b",
		"name": "John Doe",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	token.Header["kid"] = serve.LegitKid
	token.Header["x5c"] = []string{base64.StdEncoding.EncodeToString(block.Bytes)}

	return token.SignedString(key)
}

func main() {
	common.Execute(serve.RunServer, common.NewJwtCmd(generateToken))
}
