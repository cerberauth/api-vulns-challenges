package main

import (
	"github.com/cerberauth/api-vulns-challenges/challenges/tls-http-accepted/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
)

func main() {
	common.Execute(serve.RunServer)
}
