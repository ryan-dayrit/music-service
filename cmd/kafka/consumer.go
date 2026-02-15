package kafka

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"music-service/internal/config"
	"music-service/internal/handler/kafka"
	"music-service/internal/repository/postgres/orm"
	"music-service/pkg/postgres/orm/db"
)

func NewKafkaConsumerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "kafka-consumer",
		Short: "starts the kafka consumer",
		Long:  `starts the kafka consumer which listens to topics and processes messages`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			cfg, err := config.Load()
			if err != nil {
				log.Panicf("failed to load config %v", err)
			}

			db := db.NewDB(cfg.Postgres)
			defer db.Close()

			repository := orm.NewRepository(db)

			handler, err := kafka.NewConsumerHandler(cfg.Kafka, repository)
			if err != nil {
				log.Panicf("Error creating Kafka handler: %v", err)
			}

			handler.Consume(ctx)
		},
	}
}
