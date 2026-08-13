package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/cerberauth/api-vulns-challenges/challenges/jwt-facebook-token-relay/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

const AttackerAudience = "5566778899001122"

func generateToken(audience string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	privateKeyBytes, err := os.ReadFile(path.Join(cwd, "keys", "idp_private_key.pem"))
	if err != nil {
		return "", err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", err
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": serve.Issuer,
		"aud": audience,
		"sub": "10159876543210987",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}).SignedString(key)
}

func newJwtCmd(use, audience string) *cobra.Command {
	return &cobra.Command{
		Use: use,
		Run: func(cmd *cobra.Command, args []string) {
			tokenString, err := generateToken(audience)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(tokenString)
		},
	}
}

func main() {
	common.Execute(
		serve.RunServer,
		newJwtCmd("jwt-attacker-app", AttackerAudience),
		newJwtCmd("jwt-victim-app", serve.VictimAudience),
	)
}
