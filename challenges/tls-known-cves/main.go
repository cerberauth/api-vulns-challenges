package main

import (
	"github.com/cerberauth/api-vulns-challenges/challenges/tls-known-cves/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
)

func main() {
	common.Execute(serve.RunServer)
}
