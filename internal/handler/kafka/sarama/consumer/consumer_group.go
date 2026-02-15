package consumer

import (
	"log"

	"github.com/IBM/sarama"

	"music-service/internal/handler/kafka"
)

type consumerGroupHandler struct {
	Ready                 chan bool
	MessageValueProcessor *kafka.MessageValueProcessor
}

func NewConsumerGroupHandler(ready chan bool, messageValueProcessor *kafka.MessageValueProcessor) *consumerGroupHandler {
	return &consumerGroupHandler{
		Ready:                 ready,
		MessageValueProcessor: messageValueProcessor,
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

			h.MessageValueProcessor.ProcessMessageValue(message.Value)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
