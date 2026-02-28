package main

import (
	"github.com/cerberauth/api-vulns-challenges/challenges/strong-api-key/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
)

func main() {
	common.Execute(serve.RunServer)
}
