package kafka

import (
	"context"
	"log"

	"music-service/internal/config"
	handler "music-service/internal/handler/kafka"

	"github.com/spf13/cobra"
)

func NewProducerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "producer",
		Short: "starts the kafka producer",
		Long:  `starts the kafka producer which sends messages to topics`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			cfg, err := config.Load()
			if err != nil {
				log.Panicf("failed to load config %v", err)
			}

			producer, err := handler.NewProducer(cfg.Kafka)
			if err != nil {
				log.Panicf("Error creating Kafka producer: %v", err)
			}

			producer.Produce(ctx)
		},
	}
}
