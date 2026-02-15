package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"music-service/cmd/grpc"
	"music-service/cmd/kafka"
	"music-service/cmd/postgres"
	"music-service/cmd/rest"
	rest_client "music-service/cmd/rest/client"
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

	rootCmd.AddCommand(kafka.NewKafkaConsumerCommand())
	rootCmd.AddCommand(kafka.NewKafkaProducerCommand())

	rootCmd.AddCommand(postgres.NewPostgresGetAllCommand())
	rootCmd.AddCommand(postgres.NewPostgresGetByIdCommand())
	rootCmd.AddCommand(postgres.NewPostgresInsertCommand())

	rootCmd.AddCommand(rest_client.NewRestClientSingleCommand())
	rootCmd.AddCommand(rest_client.NewRestClientMultiCommand())
	rootCmd.AddCommand(rest.NewRestServerCommand())
}
