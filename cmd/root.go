package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(clientCmd)
	rootCmd.AddCommand(consumerCmd)
	rootCmd.AddCommand(producerCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(toolCmd)
}
