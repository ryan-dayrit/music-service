package cmd

import (
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "shows the albums returned from the gRPC service",
	Long:  `calls the gRPC service and shows the albums returned`,
	Run: func(cmd *cobra.Command, args []string) {
		print("client")
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
