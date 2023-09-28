package main

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Generate JWT with alg none
	// token := jwt.New(jwt.SigningMethodNone)
	// tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	// Generate JWT with weak HMAC secret
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte("secret"))

	fmt.Println(tokenString, err)
}
