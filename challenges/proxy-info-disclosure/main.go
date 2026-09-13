package main

import (
	"github.com/cerberauth/api-vulns-challenges/challenges/proxy-info-disclosure/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
)

func main() {
	common.Execute(serve.RunServer)
}
