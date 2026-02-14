package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"music-service/cmd/grpc"
	"music-service/cmd/kafka"
	"music-service/cmd/postgres"
	"music-service/cmd/rest"
)

var rootCmd = &cobra.Command{}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(grpc.NewGrpcClientCommand())
	rootCmd.AddCommand(grpc.NewGrpcServerCommand())
	rootCmd.AddCommand(kafka.NewConsumerCommand())
	rootCmd.AddCommand(kafka.NewProducerCommand())
	rootCmd.AddCommand(postgres.NewToolCommand())
	rootCmd.AddCommand(rest.NewRestClientCommand())
	rootCmd.AddCommand(rest.NewRestServerCommand())
}
