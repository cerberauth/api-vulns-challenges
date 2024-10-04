package jwt

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

func NewJwtCmd() (jwtCmd *cobra.Command) {
	jwtCmd = &cobra.Command{
		Use: "jwt",
		Run: func(cmd *cobra.Command, args []string) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub":  "2cb307ba-bb46-4194-854f-4774046d9c9b",
				"name": "John Doe",
				"iat":  time.Now().Unix(),
				"exp":  time.Now().Add(time.Hour).Unix(),
			})
			tokenString, err := token.SignedString([]byte("secret"))
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(tokenString)
		},
	}

	return jwtCmd
}
