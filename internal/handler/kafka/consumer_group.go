package kafka

import (
	"log"
	"math/rand/v2"

	"github.com/IBM/sarama"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"

	"music-service/gen/pb"
	"music-service/internal/models"
	"music-service/internal/repository/postgres/orm"
)

type consumerGroupHandler struct {
	Ready      chan bool
	Repository orm.Repository
}

func NewConsumerGroupHandler(ready chan bool, repository orm.Repository) *consumerGroupHandler {
	return &consumerGroupHandler{
		Ready:      ready,
		Repository: repository,
	}
}

func (cgh *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	close(cgh.Ready)
	return nil
}

func (_ *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			protoAlbum := &pb.Album{}
			if err := proto.Unmarshal(message.Value, protoAlbum); err != nil {
				log.Fatalf("Failed to unmarshall to album: %v", err)
			}

			if h.Repository != nil {
				existingAlbum, err := h.Repository.Read(int(protoAlbum.Id))
				if err != nil && err.Error() != "pg: no rows in result set" {
					log.Fatalf("Failed to read album from postgres: %v", err)
				}

				if existingAlbum == nil {
					newAlbum := models.Album{
						Id:     int(protoAlbum.Id),
						Title:  protoAlbum.Title,
						Artist: protoAlbum.Artist,
						Price:  decimal.NewFromFloat(rand.Float64()),
					}
					err = h.Repository.Create(newAlbum)
					if err != nil {
						log.Fatalf("Failed to create album in postgres: %v", err)
					}
				}
			}

			log.Printf("Message claimed: value = %s (Id=%d, Title=%s, Artist=%s, Price=%.2f), timestamp = %v, topic = %s", string(message.Value), protoAlbum.Id, protoAlbum.Title, protoAlbum.Artist, protoAlbum.Price, message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
