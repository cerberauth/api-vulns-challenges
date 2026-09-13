package main

import (
	"github.com/cerberauth/api-vulns-challenges/challenges/proxy-client-ip-spoofing/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
)

func main() {
	common.Execute(serve.RunServer)
}
