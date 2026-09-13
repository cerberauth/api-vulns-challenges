package main

import (
	"github.com/cerberauth/api-vulns-challenges/challenges/proxy-host-header-injection/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
)

func main() {
	common.Execute(serve.RunServer)
}
