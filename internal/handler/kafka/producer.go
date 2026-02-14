package kafka

import (
	"context"
	"log"

	"music-service/gen/pb"
	"music-service/pkg/kafka"

	"math/rand/v2"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
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
	album := &pb.Album{
		Id:     rand.Int32(),
		Title:  uuid.NewString(),
		Artist: uuid.NewString(),
		Price:  rand.Float32(),
	}
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
