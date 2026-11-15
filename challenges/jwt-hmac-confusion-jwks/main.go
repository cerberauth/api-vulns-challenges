package main

import (
	"github.com/cerberauth/api-vulns-challenges/challenges/jwt-hmac-confusion-jwks/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
)

func main() {
	common.Execute(serve.RunServer)
}
