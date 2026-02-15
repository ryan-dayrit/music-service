package sarama

import (
	"context"
	"log"
	"math/rand/v2"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"music-service/gen/pb"
	"music-service/internal/config"
	"music-service/internal/handler/kafka/sarama/producer"
)

func NewKafkaProducerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "kafka-producer-sarama",
		Short: "starts the kafka producer implemented with the sarama library",
		Long:  `starts the kafka producer which sends messages to topics using the sarama library`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			cfg, err := config.Load()
			if err != nil {
				log.Panicf("failed to load config %v", err)
			}

			producerHandler, err := producer.NewProducerHandler(cfg.Kafka)
			if err != nil {
				log.Panicf("error creating Kafka producer: %v", err)
			}

			album := &pb.Album{
				Id:     rand.Int32(),
				Title:  uuid.NewString(),
				Artist: uuid.NewString(),
				Price:  rand.Float32(),
			}
			producerHandler.Produce(ctx, album)
		},
	}
}
