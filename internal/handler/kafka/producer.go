package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"

	"music-service/gen/pb"
	"music-service/pkg/kafka"
)

type Producer interface {
	Produce(ctx context.Context, album *pb.Album)
}

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

func (p *producer) Produce(ctx context.Context, album *pb.Album) {
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
