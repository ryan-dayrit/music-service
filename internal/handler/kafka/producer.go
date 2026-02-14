package kafka

import (
	"context"
	"fmt"
	"log"

	"music-service/pkg/kafka"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type producer struct {
	cfg          kafka.Config
	syncProducer sarama.SyncProducer
}

func NewProducer(cfg kafka.Config) (*producer, error) {
	syncProducer, err := kafka.NewSyncProducer(cfg)
	if err != nil {
		return nil, err
	}
	return &producer{cfg: cfg, syncProducer: syncProducer}, nil
}

func (p *producer) Produce(ctx context.Context) {
	msgStr := fmt.Sprintf("random message - %s", uuid.NewString())
	msg := &sarama.ProducerMessage{
		Topic: p.cfg.Topics,
		Value: sarama.StringEncoder(msgStr),
	}
	partition, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		log.Panicf("Failed to send message: %v", err)
	} else {
		log.Printf("message sent; partition=%d,offset=%d", partition, offset)
	}
}
