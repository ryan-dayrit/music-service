package consumer

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"

	internal_kafka "music-service/internal/handler/kafka"
	"music-service/internal/repository/postgres/orm"
	"music-service/pkg/kafka"
	sarama_wrapper "music-service/pkg/kafka/sarama"
)

type consumerHandler struct {
	cfg           kafka.Config
	consumerGroup sarama.ConsumerGroup
	repository    orm.Repository
}

func NewConsumerHandler(cfg kafka.Config, repository orm.Repository) (kafka.ConsumerHandler, error) {
	consumerGroup, err := sarama_wrapper.NewConsumerGroup(cfg)
	if err != nil {
		return nil, err
	}

	return &consumerHandler{
		cfg:           cfg,
		consumerGroup: consumerGroup,
		repository:    repository,
	}, nil
}

func (h *consumerHandler) Consume(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	messageValueProcessor := internal_kafka.NewMessageValueProcessor(h.repository)
	consumerGroupHandler := NewConsumerGroupHandler(make(chan bool), messageValueProcessor)

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := h.consumerGroup.Consume(ctx, strings.Split(h.cfg.Topics, ","), consumerGroupHandler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumerGroupHandler.Ready = make(chan bool)
		}
	}()

	<-consumerGroupHandler.Ready

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning := true; keepRunning; {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(h.consumerGroup, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err := h.consumerGroup.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}

	return nil
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}
