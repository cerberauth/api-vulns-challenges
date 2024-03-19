package jwt

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

func GenerateRS512JWT(sub string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	privateKeyBytes, err := os.ReadFile(cwd + string(os.PathSeparator) + "keys/private_key.pem")
	if err != nil {
		return "", err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", err
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": sub,
	}).SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewJwtCmd() (jwtCmd *cobra.Command) {
	jwtCmd = &cobra.Command{
		Use: "jwt",
		Run: func(cmd *cobra.Command, args []string) {
			tokenString, err := GenerateRS512JWT("abc123")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(tokenString)
		},
	}

	return jwtCmd
}
