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
	return &consumer{
		confluentConsumer:     confluentConsumer,
		messageValueProcessor: messageValueProcessor,
		parallelWorkers:       parallelWorkers,
	}
}

func (c *consumer) Consume(ctx context.Context) error {
	tasks := make(chan *kafka.Message, 1000)
	acks := make(chan ack, 1000)

	var wg sync.WaitGroup
	for i := 0; i < c.parallelWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			log.Printf("worker %d started", workerID)

			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-tasks:
					if !ok {
						return
					}

					log.Printf(
						"Processing message from %s[%d]@%d",
						*msg.TopicPartition.Topic,
						msg.TopicPartition.Partition,
						msg.TopicPartition.Offset,
					)

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
		ticker := time.NewTicker(5 * time.Second)
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
					offsets, err := c.confluentConsumer.Commit()
					if err != nil {
						log.Printf("Commit failed: %v", err)
					} else {
						log.Printf("Committed %d partitions", len(offsets))
					}
					pending = make(map[kafka.TopicPartition]kafka.Offset)
				}
			}
		}
	}()

	log.Println("starting consumer polling loop")

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down consumer...")
			close(tasks)
			wg.Wait()

			_, err := c.confluentConsumer.Commit()
			if err != nil {
				log.Printf("final commit failed: %v", err)
			}

			return c.confluentConsumer.Close()

		default:
			ev := c.confluentConsumer.Poll(100)
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
