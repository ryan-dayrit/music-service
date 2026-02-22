package producer

import (
	"context"
	"fmt"
	"log"
	"strings"

	ext_kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"google.golang.org/protobuf/proto"

	"music-service/gen/pb"
	"music-service/pkg/kafka"
)

type producerHandler struct {
	cfg               kafka.Config
	topic             string
	confluentProducer *ext_kafka.Producer
}

func NewProducerHandler(cfg kafka.Config) (kafka.ProducerHandler, error) {
	extCfg := &ext_kafka.ConfigMap{"bootstrap.servers": cfg.Brokers}

	confluentProducer, err := ext_kafka.NewProducer(extCfg)
	if err != nil {
		return nil, err
	}

	// Use primary topic for produce (first topic if comma-separated)
	rawTopic, _, _ := strings.Cut(cfg.Topics, ",")
	topic := strings.TrimSpace(rawTopic)
	if topic == "" {
		confluentProducer.Close()
		return nil, fmt.Errorf("no topics configured")
	}

	p := &producerHandler{
		cfg:               cfg,
		topic:             topic,
		confluentProducer: confluentProducer,
	}

	go p.runDeliveryReports()

	return p, nil
}

func (p *producerHandler) runDeliveryReports() {
	for e := range p.confluentProducer.Events() {
		msg, ok := e.(*ext_kafka.Message)
		if !ok {
			continue
		}
		if msg.TopicPartition.Error != nil {
			log.Printf("failed to deliver message: %v", msg.TopicPartition.Error)
		}
	}
}

func (p *producerHandler) Produce(ctx context.Context, album *pb.Album) {
	select {
	case <-ctx.Done():
		log.Printf("produce cancelled: %v", ctx.Err())
		return
	default:
	}

	marshaledAlbum, err := proto.Marshal(album)
	if err != nil {
		log.Panicf("failed to marshal album: %v", err)
	}

	err = p.confluentProducer.Produce(&ext_kafka.Message{
		TopicPartition: ext_kafka.TopicPartition{Topic: &p.topic, Partition: ext_kafka.PartitionAny},
		Value:          marshaledAlbum,
	}, nil)
	if err != nil {
		log.Panicf("failed to produce album: %v", err)
	}
}
