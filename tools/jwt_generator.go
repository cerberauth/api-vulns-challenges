package main

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	token := jwt.New(jwt.SigningMethodNone)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	fmt.Println(tokenString, err)
}
