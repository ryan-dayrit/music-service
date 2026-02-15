package producer

import (
	"context"
	"log"

	ext_kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"google.golang.org/protobuf/proto"

	"music-service/gen/pb"
	"music-service/pkg/kafka"
)

type producerHandler struct {
	cfg               kafka.Config
	confluentProducer *ext_kafka.Producer
}

func NewProducerHandler(cfg kafka.Config) (kafka.ProducerHandler, error) {
	extCfg := &ext_kafka.ConfigMap{"bootstrap.servers": cfg.Brokers}

	confluentProducer, err := ext_kafka.NewProducer(extCfg)
	if err != nil {
		return nil, err
	}

	return &producerHandler{cfg: cfg, confluentProducer: confluentProducer}, nil
}

func (p *producerHandler) Produce(ctx context.Context, album *pb.Album) {
	deliveryChan := make(chan ext_kafka.Event)

	marshaledAlbum, err := proto.Marshal(album)
	if err != nil {
		log.Panicf("failed to marshal album: %v", err)
	}

	err = p.confluentProducer.Produce(&ext_kafka.Message{
		TopicPartition: ext_kafka.TopicPartition{Topic: &p.cfg.Topics, Partition: ext_kafka.PartitionAny},
		Value:          marshaledAlbum,
		Headers:        []ext_kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, deliveryChan)
	if err != nil {
		log.Panicf("failed to produce album: %v", err)
	}

	event := <-deliveryChan
	message := event.(*ext_kafka.Message)

	if message.TopicPartition.Error != nil {
		log.Printf("failed to deliver message: %v\n", message.TopicPartition.Error)
	} else {
		log.Printf("delivered message to topic %s [%d] at offset %v\n",
			*message.TopicPartition.Topic, message.TopicPartition.Partition, message.TopicPartition.Offset)
	}

	close(deliveryChan)
}
