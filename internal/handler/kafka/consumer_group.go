package kafka

import (
	"log"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"

	"music-service/gen/pb"
)

type consumerGroupHandler struct {
	Ready chan bool
}

func (cgh *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	close(cgh.Ready)
	return nil
}

func (_ *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (_ *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			album := &pb.Album{}
			if err := proto.Unmarshal(message.Value, album); err != nil {
				log.Fatalf("Failed to unmarshall to album: %v", err)
			}

			log.Printf("Message claimed: value = %s (Id=%d, Title=%s, Artist=%s, Price=%.2f), timestamp = %v, topic = %s", string(message.Value), album.Id, album.Title, album.Artist, album.Price, message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
