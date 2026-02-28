package main

import (
	"github.com/cerberauth/api-vulns-challenges/challenges/strong-http-basic/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
)

func main() {
	common.Execute(serve.RunServer)
}
