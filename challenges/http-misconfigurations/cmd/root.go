package cmd

import (
	"os"

	"github.com/cerberauth/api-vulns-challenges/challenges/http-misconfigurations/cmd/serve"

	"github.com/spf13/cobra"
)

func NewRootCmd() (cmd *cobra.Command) {
	var rootCmd = &cobra.Command{}

	rootCmd.AddCommand(serve.NewServeCmd())

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	c := NewRootCmd()

	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
