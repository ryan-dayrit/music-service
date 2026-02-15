package producer

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"

	"music-service/gen/pb"
	"music-service/pkg/kafka"
	sarama_wrapper "music-service/pkg/kafka/sarama"
)

type producerHandler struct {
	cfg          kafka.Config
	syncProducer sarama.SyncProducer
}

func NewProducerHandler(cfg kafka.Config) (kafka.ProducerHandler, error) {
	syncProducer, err := sarama_wrapper.NewSyncProducer(cfg)
	if err != nil {
		return nil, err
	}
	return &producerHandler{cfg: cfg, syncProducer: syncProducer}, nil
}

func (p *producerHandler) Produce(ctx context.Context, album *pb.Album) {
	marshaledAlbum, err := proto.Marshal(album)
	if err != nil {
		log.Panicf("Failed to marshal album: %v", err)
	}
	msg := &sarama.ProducerMessage{
		Topic: p.cfg.Topics,
		Value: sarama.ByteEncoder(marshaledAlbum),
	}
	partition, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		log.Panicf("Failed to send message: %v", err)
	} else {
		log.Printf("message sent (Id=%d, Title=%s, Artist=%s, Price=%.2f); partition=%d,offset=%d", album.Id, album.Title, album.Artist, album.Price, partition, offset)
	}
}
