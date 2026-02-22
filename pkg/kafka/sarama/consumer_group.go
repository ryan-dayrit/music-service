package sarama

import (
	"fmt"
	"strings"

	"github.com/IBM/sarama"

	"music-service/pkg/kafka"
)

func NewConsumerGroup(cfg kafka.Config) (sarama.ConsumerGroup, error) {
	saramaCfg := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(sarama.DefaultVersion.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse Kafka version: %w", err)
	}
	saramaCfg.Version = version

	brokers := strings.Split(cfg.Brokers, ",")
	for i := range brokers {
		brokers[i] = strings.TrimSpace(brokers[i])
	}

	switch cfg.Assignor {
	case "sticky":
		saramaCfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		saramaCfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		saramaCfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		return nil, fmt.Errorf("Unrecognized consumer group partition assignor: %s", cfg.Assignor)
	}

	if cfg.Oldest {
		saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	consumerGroup, err := sarama.NewConsumerGroup(brokers, cfg.ConsumerGroup, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("Error creating consumer group: %v", err)
	}
	return consumerGroup, nil
}
