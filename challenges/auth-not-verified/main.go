package main

import (
	"github.com/cerberauth/api-vulns-challenges/challenges/auth-not-verified/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
)

func main() {
	common.Execute(serve.RunServer)
}
