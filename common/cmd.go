package common

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func NewJwtCmd(generate func(vulnerable bool) (string, error)) *cobra.Command {
	var vulnerable bool
	cmd := &cobra.Command{
		Use: "jwt",
		Run: func(cmd *cobra.Command, args []string) {
			tokenString, err := generate(vulnerable)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(tokenString)
		},
	}
	cmd.Flags().BoolVar(&vulnerable, "vulnerable", true, "Generate a token as issued in the vulnerable setup. Set to false to generate a token matching the fixed, non-vulnerable server")
	return cmd
}

func Execute(runner func(port string, vulnerable bool), additionalCmds ...*cobra.Command) {
	var port string
	var vulnerable bool
	serveCmd := &cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) { runner(port, vulnerable) },
	}
	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to listen on")
	serveCmd.Flags().BoolVar(&vulnerable, "vulnerable", true, "Run the server in its vulnerable mode. Set to false to run the fixed, non-vulnerable implementation")

	rootCmd := &cobra.Command{Use: "app"}
	rootCmd.AddCommand(serveCmd)
	for _, cmd := range additionalCmds {
		rootCmd.AddCommand(cmd)
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
