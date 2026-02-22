package producer

import (
	"context"
	"log"
	"strings"
	"sync"

	ext_kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"google.golang.org/protobuf/proto"

	"music-service/gen/pb"
	"music-service/pkg/kafka"
)

const (
	deliveryChanSize = 1000
	flushTimeoutMs   = 5000
)

type producerHandler struct {
	cfg               kafka.Config
	topic             string
	confluentProducer *ext_kafka.Producer
	deliveryChan      chan ext_kafka.Event
	wg                sync.WaitGroup
}

func NewProducerHandler(cfg kafka.Config) (kafka.ProducerHandler, error) {
	extCfg := &ext_kafka.ConfigMap{"bootstrap.servers": cfg.Brokers}

	confluentProducer, err := ext_kafka.NewProducer(extCfg)
	if err != nil {
		return nil, err
	}

	// Use primary topic for produce (first topic if comma-separated)
	topic, _, _ := strings.Cut(cfg.Topics, ",")

	p := &producerHandler{
		cfg:               cfg,
		topic:             topic,
		confluentProducer: confluentProducer,
		deliveryChan:      make(chan ext_kafka.Event, deliveryChanSize),
	}

	p.wg.Add(1)
	go p.runDeliveryReports()

	return p, nil
}

func (p *producerHandler) runDeliveryReports() {
	defer p.wg.Done()
	for e := range p.deliveryChan {
		msg, ok := e.(*ext_kafka.Message)
		if !ok {
			continue
		}
		if msg.TopicPartition.Error != nil {
			log.Printf("failed to deliver message: %v", msg.TopicPartition.Error)
		}
	}
}

func (p *producerHandler) Close() {
	p.confluentProducer.Flush(flushTimeoutMs)
	close(p.deliveryChan)
	p.wg.Wait()
	p.confluentProducer.Close()
}

func (p *producerHandler) Produce(ctx context.Context, album *pb.Album) {
	marshaledAlbum, err := proto.Marshal(album)
	if err != nil {
		log.Panicf("failed to marshal album: %v", err)
	}

	err = p.confluentProducer.Produce(&ext_kafka.Message{
		TopicPartition: ext_kafka.TopicPartition{Topic: &p.topic, Partition: ext_kafka.PartitionAny},
		Value:          marshaledAlbum,
	}, p.deliveryChan)
	if err != nil {
		log.Panicf("failed to produce album: %v", err)
	}
}
