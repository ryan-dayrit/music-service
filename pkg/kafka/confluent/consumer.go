package confluent

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"music-service/pkg/kafka/message"
)

const (
	defaultTasksBufSize    = 1000
	defaultAcksBufSize     = 1000
	defaultPollTimeoutMs   = 100
	defaultCommitInterval  = 5 * time.Second
	minParallelWorkers     = 1
)

type ack struct {
	tp  kafka.TopicPartition
	off kafka.Offset
}

type consumer struct {
	confluentConsumer     *kafka.Consumer
	messageValueProcessor message.MessageValueProcessor
	parallelWorkers       int
}

func NewConsumer(confluentConsumer *kafka.Consumer, messageValueProcessor message.MessageValueProcessor, parallelWorkers int) *consumer {
	if parallelWorkers < minParallelWorkers {
		parallelWorkers = minParallelWorkers
	}
	return &consumer{
		confluentConsumer:     confluentConsumer,
		messageValueProcessor: messageValueProcessor,
		parallelWorkers:       parallelWorkers,
	}
}

func (c *consumer) Consume(ctx context.Context) error {
	tasks := make(chan *kafka.Message, defaultTasksBufSize)
	acks := make(chan ack, defaultAcksBufSize)

	var wg sync.WaitGroup
	for i := 0; i < c.parallelWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-tasks:
					if !ok {
						return
					}

					c.messageValueProcessor.Process(msg.Value)

					acks <- ack{
						tp:  msg.TopicPartition,
						off: msg.TopicPartition.Offset,
					}
				}
			}
		}(i)
	}

	go func() {
		ticker := time.NewTicker(defaultCommitInterval)
		defer ticker.Stop()

		pending := make(map[kafka.TopicPartition]kafka.Offset)

		for {
			select {
			case <-ctx.Done():
				return
			case ack := <-acks:
				nextOffset := ack.off + 1
				tp := kafka.TopicPartition{
					Topic:     ack.tp.Topic,
					Partition: ack.tp.Partition,
					Offset:    nextOffset,
				}

				pending[ack.tp] = nextOffset

				_, err := c.confluentConsumer.StoreOffsets([]kafka.TopicPartition{tp})
				if err != nil {
					log.Printf("Failed to store offset: %v", err)
				}

			case <-ticker.C:
				if len(pending) > 0 {
					_, err := c.confluentConsumer.Commit()
					if err != nil {
						log.Printf("Commit failed: %v", err)
					}
					clear(pending)
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			close(tasks)
			wg.Wait()

			_, err := c.confluentConsumer.Commit()
			if err != nil {
				log.Printf("final commit failed: %v", err)
			}

			return c.confluentConsumer.Close()

		default:
			ev := c.confluentConsumer.Poll(defaultPollTimeoutMs)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				select {
				case <-ctx.Done():
					return ctx.Err()
				case tasks <- e:
				}

			case kafka.Error:
				log.Printf("consumer error: %v", e)
				if e.Code() == kafka.ErrAllBrokersDown {
					return fmt.Errorf("all brokers down")
				}

			case kafka.AssignedPartitions:
				log.Printf("partitions assigned: %v", e)
				err := c.confluentConsumer.Assign(e.Partitions)
				if err != nil {
					log.Printf("failed to assign partitions: %v", err)
				}

			case kafka.RevokedPartitions:
				log.Printf("partitions revoked: %v", e)
				err := c.confluentConsumer.Unassign()
				if err != nil {
					log.Printf("failed to unassign partitions: %v", err)
				}

			default:
				log.Printf("ignored event: %v", e)
			}
		}
	}
}
