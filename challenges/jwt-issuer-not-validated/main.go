package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cerberauth/api-vulns-challenges/challenges/jwt-issuer-not-validated/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

func issueToken(issuer string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": issuer,
		"sub": "2cb307ba-bb46-4194-854f-4774046d9c9b",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString([]byte(serve.IdPSecret))
}

func main() {
	otherIssuerCmd := &cobra.Command{
		Use:   "other-issuer-jwt",
		Short: "Generate a token issued by another tenant of the same identity provider",
		Run: func(cmd *cobra.Command, args []string) {
			tokenString, err := issueToken(serve.OtherIssuer)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(tokenString)
		},
	}

	common.Execute(serve.RunServer, common.NewJwtCmd(func(_ bool) (string, error) {
		return issueToken(serve.Issuer)
	}), otherIssuerCmd)
}
