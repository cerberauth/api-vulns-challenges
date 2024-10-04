package jwt

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

func GenerateRS512JWT(sub string) (string, error) {
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
		"sub":  sub,
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

func NewJwtCmd() (jwtCmd *cobra.Command) {
	jwtCmd = &cobra.Command{
		Use: "jwt",
		Run: func(cmd *cobra.Command, args []string) {
			tokenString, err := GenerateRS512JWT("2cb307ba-bb46-4194-854f-4774046d9c9b")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(tokenString)
		},
	}

	return jwtCmd
}
