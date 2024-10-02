package serve

import (
	"github.com/spf13/cobra"

	"github.com/cerberauth/api-vulns-challenges/challenges/jwt-null-signature/serve"
)

var (
	port string
)

func NewServeCmd() (serveCmd *cobra.Command) {
	serveCmd = &cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) {
			serve.RunServer(port)
		},
	}

	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to listen on")

	return serveCmd
}
