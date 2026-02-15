package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"music-service/cmd/grpc"
	"music-service/cmd/kafka/confluent"
	"music-service/cmd/kafka/sarama"
	"music-service/cmd/postgres"
	rest_client "music-service/cmd/rest/client"
	rest_server "music-service/cmd/rest/server"
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

	rootCmd.AddCommand(confluent.NewKafkaConsumerCommand())

	rootCmd.AddCommand(sarama.NewKafkaConsumerCommand())
	rootCmd.AddCommand(sarama.NewKafkaProducerCommand())

	rootCmd.AddCommand(postgres.NewPostgresGetAllCommand())
	rootCmd.AddCommand(postgres.NewPostgresGetByIdCommand())
	rootCmd.AddCommand(postgres.NewPostgresInsertCommand())

	rootCmd.AddCommand(rest_client.NewRestClientSingleCommand())
	rootCmd.AddCommand(rest_client.NewRestClientMultiCommand())
	rootCmd.AddCommand(rest_server.NewRestServerCommand())
}
