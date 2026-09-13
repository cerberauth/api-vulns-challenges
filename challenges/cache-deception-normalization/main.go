package main

import (
	"github.com/cerberauth/api-vulns-challenges/challenges/cache-deception-normalization/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
)

func main() {
	common.Execute(serve.RunServer)
}
