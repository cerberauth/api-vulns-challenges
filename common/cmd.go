package common

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func NewJwtCmd(generate func() (string, error)) *cobra.Command {
	return &cobra.Command{
		Use: "jwt",
		Run: func(cmd *cobra.Command, args []string) {
			tokenString, err := generate()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(tokenString)
		},
	}
}

func Execute(runner func(port string), additionalCmds ...*cobra.Command) {
	var port string
	serveCmd := &cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) { runner(port) },
	}
	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to listen on")

	rootCmd := &cobra.Command{Use: "app"}
	rootCmd.AddCommand(serveCmd)
	for _, cmd := range additionalCmds {
		rootCmd.AddCommand(cmd)
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
