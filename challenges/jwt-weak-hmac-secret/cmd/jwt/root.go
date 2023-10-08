package jwt

import (
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

func NewJwtCmd() (jwtCmd *cobra.Command) {
	jwtCmd = &cobra.Command{
		Use: "jwt",
		Run: func(cmd *cobra.Command, args []string) {
			token := jwt.New(jwt.SigningMethodHS256)
			tokenString, err := token.SignedString([]byte("secret"))
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("token: %s", tokenString)
		},
	}

	return jwtCmd
}
