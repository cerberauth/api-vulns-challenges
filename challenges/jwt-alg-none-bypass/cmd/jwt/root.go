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
			token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
				"sub":  "2cb307ba-bb46-4194-854f-4774046d9c9b",
				"name": "John Doe",
				"iat":  1516239022,
				"exp":  1516242622,
			})
			tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(tokenString)
		},
	}

	return jwtCmd
}
