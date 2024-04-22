package serve

import (
	"github.com/spf13/cobra"

	"github.com/cerberauth/api-vulns-challenges/challenges/auth-not-verified/serve"
)

func NewServeCmd() (serveCmd *cobra.Command) {
	serveCmd = &cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) {
			serve.RunServer()
		},
	}

	return serveCmd
}
