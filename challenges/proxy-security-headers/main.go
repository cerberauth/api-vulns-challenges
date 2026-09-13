package main

import (
	"github.com/cerberauth/api-vulns-challenges/challenges/proxy-security-headers/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
)

func main() {
	common.Execute(serve.RunServer)
}
