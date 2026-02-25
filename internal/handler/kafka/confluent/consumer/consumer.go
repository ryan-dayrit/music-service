package consumer

import (
	"fmt"
	"strings"

	ext_kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"music-service/internal/handler/kafka/message"
	"music-service/internal/repository/postgres/orm"
	"music-service/pkg/kafka"
	"music-service/pkg/kafka/confluent"
)

const defaultParallelWorkers = 5
const defaultAssignor = "cooperative-sticky"

func newConsumerConfig(cfg kafka.Config) *ext_kafka.ConfigMap {
	assignor := cfg.Assignor
	if assignor == "" {
		assignor = defaultAssignor
	}
	return &ext_kafka.ConfigMap{
		"bootstrap.servers":             cfg.Brokers,
		"group.id":                      cfg.ConsumerGroup,
		"auto.offset.reset":             "earliest",
		"enable.auto.commit":            false,
		"enable.auto.offset.store":      false,
		"partition.assignment.strategy": assignor,
		"session.timeout.ms":            30000,
		"max.poll.interval.ms":          300000,
	}
}

func parseTopics(topics string) []string {
	raw := strings.Split(topics, ",")
	result := make([]string, 0, len(raw))
	for _, t := range raw {
		if trimmed := strings.TrimSpace(t); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func NewConsumerHandler(cfg kafka.Config, repository orm.Repository) (kafka.ConsumerHandler, error) {
	confluentConsumer, err := ext_kafka.NewConsumer(newConsumerConfig(cfg))
	if err != nil {
		return nil, err
	}

	topics := parseTopics(cfg.Topics)
	if len(topics) == 0 {
		_ = confluentConsumer.Close()
		return nil, fmt.Errorf("no topics configured")
	}

	err = confluentConsumer.SubscribeTopics(topics, nil)
	if err != nil {
		_ = confluentConsumer.Close()
		return nil, err
	}

	processor := message.NewMessageValueProcessor(repository)
	return confluent.NewConsumer(confluentConsumer, processor, defaultParallelWorkers), nil
}
