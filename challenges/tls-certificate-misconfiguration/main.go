package main

import (
	"log"
	"os"

	"github.com/cerberauth/api-vulns-challenges/challenges/tls-certificate-misconfiguration/serve"
	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/spf13/cobra"
)

func main() {
	caCmd := &cobra.Command{
		Use:   "ca",
		Short: "Print the CA certificate that signs the fixed server certificate",
		Run: func(cmd *cobra.Command, args []string) {
			caPEM, err := serve.CAPEM()
			if err != nil {
				log.Fatal(err)
			}
			os.Stdout.Write(caPEM)
		},
	}

	common.Execute(serve.RunServer, caCmd)
}
