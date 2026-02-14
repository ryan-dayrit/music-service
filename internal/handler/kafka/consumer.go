package kafka

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

	"music-service/pkg/kafka"
)

type consumerHandler struct {
	cfg           kafka.Config
	consumerGroup sarama.ConsumerGroup
}

func NewConsumerHandler(cfg kafka.Config) (*consumerHandler, error) {
	consumerGroup, err := kafka.NewConsumerGroup(cfg)
	if err != nil {
		return nil, err
	}
	return &consumerHandler{
		cfg:           cfg,
		consumerGroup: consumerGroup,
	}, nil
}

func (h *consumerHandler) Consume(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	consumerGroupHandler := consumerGroupHandler{
		Ready: make(chan bool),
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := h.consumerGroup.Consume(ctx, strings.Split(h.cfg.Topics, ","), &consumerGroupHandler); err != nil {
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
