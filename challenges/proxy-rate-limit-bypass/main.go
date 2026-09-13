package main

import (
	"github.com/cerberauth/api-vulns-challenges/challenges/proxy-rate-limit-bypass/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
)

func main() {
	common.Execute(serve.RunServer)
}
