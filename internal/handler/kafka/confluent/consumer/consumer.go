package consumer

import (
	"strings"

	ext_kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	internal_kafka "music-service/internal/handler/kafka"
	"music-service/internal/repository/postgres/orm"
	"music-service/pkg/kafka"
	"music-service/pkg/kafka/confluent"
)

type consumerHandler struct {
	consumer kafka.ConsumerHandler
}

func NewConsumerHandler(cfg kafka.Config, repository orm.Repository) (kafka.ConsumerHandler, error) {
	extCfg := &ext_kafka.ConfigMap{
		"bootstrap.servers":             cfg.Brokers,
		"group.id":                      cfg.ConsumerGroup,
		"auto.offset.reset":             "earliest",
		"enable.auto.commit":            false,
		"enable.auto.offset.store":      false,
		"partition.assignment.strategy": "cooperative-sticky",
		"session.timeout.ms":            30000,
		"max.poll.interval.ms":          300000,
	}

	confluentConsumer, err := ext_kafka.NewConsumer(extCfg)
	if err != nil {
		return nil, err
	}

	err = confluentConsumer.SubscribeTopics(strings.Split(cfg.Topics, ","), nil)
	if err != nil {
		return nil, err
	}

	messageValueProcessor := internal_kafka.NewMessageValueProcessor(repository)

	consumerHandler := confluent.NewConsumer(confluentConsumer, messageValueProcessor, 5)
	return consumerHandler, nil
}
