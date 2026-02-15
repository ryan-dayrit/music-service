package sarama

import (
	"fmt"
	"strings"

	"music-service/pkg/kafka"

	"github.com/IBM/sarama"
)

func NewSyncProducer(cfg kafka.Config) (sarama.SyncProducer, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Version, _ = sarama.ParseKafkaVersion(sarama.DefaultVersion.String())
	saramaCfg.Producer.RequiredAcks = sarama.WaitForAll
	saramaCfg.Producer.Retry.Max = 10
	saramaCfg.Producer.Return.Successes = true

	syncProducer, err := sarama.NewSyncProducer(strings.Split(cfg.Brokers, ","), saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("Error creating consumer group: %v", err)
	}

	return syncProducer, nil
}
