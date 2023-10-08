package cmd

import (
	"os"

	"github.com/cerberauth/vulns-challenges/challenges/jwt-not-verified/cmd/jwt"
	"github.com/cerberauth/vulns-challenges/challenges/jwt-not-verified/cmd/serve"

	"github.com/spf13/cobra"
)

func NewRootCmd() (cmd *cobra.Command) {
	var rootCmd = &cobra.Command{}

	rootCmd.AddCommand(serve.NewServeCmd())
	rootCmd.AddCommand(jwt.NewJwtCmd())

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
