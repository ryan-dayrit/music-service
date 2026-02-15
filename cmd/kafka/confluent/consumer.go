package confluent

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"music-service/internal/config"
	"music-service/internal/handler/kafka/confluent/consumer"
	"music-service/internal/repository/postgres/orm"
	"music-service/pkg/postgres/orm/db"
)

func NewKafkaConsumerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "kafka-consumer-confluent",
		Short: "starts the kafka consumer implemented with the confluent library",
		Long:  `starts the kafka consumer which listens to topics and processes messages using the confluent library`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			cfg, err := config.Load()
			if err != nil {
				log.Panicf("failed to load config %v", err)
			}

			db := db.NewDB(cfg.Postgres)
			defer db.Close()

			repository := orm.NewRepository(db)

			handler, err := consumer.NewConsumerHandler(cfg.Kafka, repository)
			if err != nil {
				log.Panicf("error creating consumer handler: %v", err)
			}

			handler.Consume(ctx)
		},
	}
}
