package kafka

import (
	"fmt"
	"strings"

	"github.com/IBM/sarama"
)

func NewConsumerGroup(cfg Config) (sarama.ConsumerGroup, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Version, _ = sarama.ParseKafkaVersion(sarama.DefaultVersion.String())

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

	consumerGroup, err := sarama.NewConsumerGroup(strings.Split(cfg.Brokers, ","), cfg.ConsumerGroup, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("Error creating consumer group: %v", err)
	}
	return consumerGroup, nil
}
