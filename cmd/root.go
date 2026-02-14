package cmd

import (
	"os"

	"music-service/cmd/grpc"
	"music-service/cmd/kafka"
	"music-service/cmd/postgres"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(grpc.NewClientCommand())
	rootCmd.AddCommand(grpc.NewServerCommand())
	rootCmd.AddCommand(kafka.NewConsumerCommand())
	rootCmd.AddCommand(kafka.NewProducerCommand())
	rootCmd.AddCommand(postgres.NewToolCommand())
}
