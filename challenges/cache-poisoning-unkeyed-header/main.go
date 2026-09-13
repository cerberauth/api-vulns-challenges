package main

import (
	"github.com/cerberauth/api-vulns-challenges/challenges/cache-poisoning-unkeyed-header/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
)

func main() {
	common.Execute(serve.RunServer)
}
