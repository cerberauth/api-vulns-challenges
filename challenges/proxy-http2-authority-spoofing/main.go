package main

import (
	"github.com/cerberauth/api-vulns-challenges/challenges/proxy-http2-authority-spoofing/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
)

func main() {
	common.Execute(serve.RunServer)
}
