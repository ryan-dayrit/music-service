package cmd

import (
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts the gRPC service",
	Long:  `starts the gRPC service which returns albums`,
	Run: func(cmd *cobra.Command, args []string) {
		print("server")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
