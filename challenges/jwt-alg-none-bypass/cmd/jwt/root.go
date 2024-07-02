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
			claims := &jwt.RegisteredClaims{
				Subject:  "1234567890",
				IssuedAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
			tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(tokenString)
		},
	}

	return jwtCmd
}
