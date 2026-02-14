package cmd

import (
	"context"
	"log"

	"music-service/internal/config"
	"music-service/internal/handler/kafka"

	"github.com/spf13/cobra"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "starts the kafka consumer",
	Long:  `starts the kafka consumer which listens to topics and processes messages`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		cfg, err := config.Load()
		if err != nil {
			log.Panicf("failed to load config %v", err)
		}

		handler, err := kafka.NewConsumerHandler(cfg.Kafka)
		if err != nil {
			log.Panicf("Error creating Kafka handler: %v", err)
		}

		handler.Consume(ctx)
	},
}
