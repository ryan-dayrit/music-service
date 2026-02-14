package kafka

import (
	"context"
	"log"

	"math/rand/v2"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"music-service/gen/pb"
	"music-service/internal/config"
	"music-service/internal/handler/kafka"
)

func NewProducerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "kafka-producer",
		Short: "starts the kafka producer",
		Long:  `starts the kafka producer which sends messages to topics`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			cfg, err := config.Load()
			if err != nil {
				log.Panicf("failed to load config %v", err)
			}

			producer, err := kafka.NewProducer(cfg.Kafka)
			if err != nil {
				log.Panicf("Error creating Kafka producer: %v", err)
			}

			album := &pb.Album{
				Id:     rand.Int32(),
				Title:  uuid.NewString(),
				Artist: uuid.NewString(),
				Price:  rand.Float32(),
			}
			producer.Produce(ctx, album)
		},
	}
}
