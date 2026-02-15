package consumer

import (
	"context"
	"testing"
	"time"

	"github.com/IBM/sarama"

	"music-service/internal/handler/kafka"
)

// MockConsumerGroupSession is a mock implementation of sarama.ConsumerGroupSession
type MockConsumerGroupSession struct {
	markedMessages []*sarama.ConsumerMessage
	ctx            context.Context
	cancel         context.CancelFunc
}

func NewMockConsumerGroupSession() *MockConsumerGroupSession {
	ctx, cancel := context.WithCancel(context.Background())
	return &MockConsumerGroupSession{
		markedMessages: make([]*sarama.ConsumerMessage, 0),
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (m *MockConsumerGroupSession) Claims() map[string][]int32 {
	return map[string][]int32{"test-topic": {0}}
}

func (m *MockConsumerGroupSession) MemberID() string {
	return "test-member"
}

func (m *MockConsumerGroupSession) GenerationID() int32 {
	return 1
}

func (m *MockConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
}

func (m *MockConsumerGroupSession) Commit() {
}

func (m *MockConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
}

func (m *MockConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	m.markedMessages = append(m.markedMessages, msg)
}

func (m *MockConsumerGroupSession) Context() context.Context {
	return m.ctx
}

// MockConsumerGroupClaim is a mock implementation of sarama.ConsumerGroupClaim
type MockConsumerGroupClaim struct {
	messages chan *sarama.ConsumerMessage
}

func NewMockConsumerGroupClaim() *MockConsumerGroupClaim {
	return &MockConsumerGroupClaim{
		messages: make(chan *sarama.ConsumerMessage, 10),
	}
}

func (m *MockConsumerGroupClaim) Topic() string {
	return "test-topic"
}

func (m *MockConsumerGroupClaim) Partition() int32 {
	return 0
}

func (m *MockConsumerGroupClaim) InitialOffset() int64 {
	return 0
}

func (m *MockConsumerGroupClaim) HighWaterMarkOffset() int64 {
	return 100
}

func (m *MockConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	return m.messages
}

// MockMessageValueProcessor is a mock implementation of kafka.MessageValueProcessor
type MockMessageValueProcessor struct {
	processedMessages [][]byte
}

func (m *MockMessageValueProcessor) ProcessMessageValue(value []byte) {
	m.processedMessages = append(m.processedMessages, value)
}

func TestNewConsumerGroupHandler(t *testing.T) {
	t.Run("creates handler successfully", func(t *testing.T) {
		ready := make(chan bool)
		processor := &kafka.MessageValueProcessor{}

		handler := NewConsumerGroupHandler(ready, processor)

		if handler == nil {
			t.Error("Expected handler to be created, got nil")
		}

		if handler.Ready != ready {
			t.Error("Expected Ready channel to match provided channel")
		}

		if handler.MessageValueProcessor != processor {
			t.Error("Expected MessageValueProcessor to match provided processor")
		}
	})

	t.Run("creates handler with nil processor", func(t *testing.T) {
		ready := make(chan bool)

		handler := NewConsumerGroupHandler(ready, nil)

		if handler == nil {
			t.Error("Expected handler to be created even with nil processor")
		}

		if handler.MessageValueProcessor != nil {
			t.Error("Expected MessageValueProcessor to be nil")
		}
	})
}

func TestConsumerGroupHandler_Setup(t *testing.T) {
	t.Run("closes ready channel on setup", func(t *testing.T) {
		ready := make(chan bool)
		processor := &kafka.MessageValueProcessor{}
		handler := NewConsumerGroupHandler(ready, processor)

		session := NewMockConsumerGroupSession()

		go func() {
			err := handler.Setup(session)
			if err != nil {
				t.Errorf("Setup should not return error, got: %v", err)
			}
		}()

		select {
		case <-ready:
			// Success - channel was closed
		case <-time.After(1 * time.Second):
			t.Error("Ready channel was not closed within timeout")
		}
	})

	t.Run("setup returns no error", func(t *testing.T) {
		ready := make(chan bool)
		processor := &kafka.MessageValueProcessor{}
		handler := NewConsumerGroupHandler(ready, processor)

		session := NewMockConsumerGroupSession()

		err := handler.Setup(session)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})
}

func TestConsumerGroupHandler_Cleanup(t *testing.T) {
	t.Run("cleanup returns no error", func(t *testing.T) {
		ready := make(chan bool)
		processor := &kafka.MessageValueProcessor{}
		handler := NewConsumerGroupHandler(ready, processor)

		session := NewMockConsumerGroupSession()

		err := handler.Cleanup(session)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})
}

func TestConsumerGroupHandler_ConsumeClaim(t *testing.T) {
	t.Run("stops consuming when context is done", func(t *testing.T) {
		ready := make(chan bool)
		processor := &kafka.MessageValueProcessor{}
		handler := NewConsumerGroupHandler(ready, processor)

		session := NewMockConsumerGroupSession()
		claim := NewMockConsumerGroupClaim()

		// Cancel context immediately
		session.cancel()

		err := handler.ConsumeClaim(session, claim)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("handles closed message channel", func(t *testing.T) {
		ready := make(chan bool)
		processor := &kafka.MessageValueProcessor{}
		handler := NewConsumerGroupHandler(ready, processor)

		session := NewMockConsumerGroupSession()
		claim := NewMockConsumerGroupClaim()

		// Close messages channel immediately
		close(claim.messages)

		err := handler.ConsumeClaim(session, claim)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})
}

func TestConsumerGroupHandler_StructFields(t *testing.T) {
	t.Run("handler has required fields", func(t *testing.T) {
		ready := make(chan bool)
		processor := &kafka.MessageValueProcessor{}
		handler := NewConsumerGroupHandler(ready, processor)

		if handler.Ready == nil {
			t.Error("Expected Ready field to be set")
		}

		if handler.MessageValueProcessor == nil {
			t.Error("Expected MessageValueProcessor field to be set")
		}
	})
}

func TestConsumerGroupHandler_InterfaceCompliance(t *testing.T) {
	t.Run("implements sarama.ConsumerGroupHandler interface", func(t *testing.T) {
		ready := make(chan bool)
		processor := &kafka.MessageValueProcessor{}
		handler := NewConsumerGroupHandler(ready, processor)

		var _ sarama.ConsumerGroupHandler = handler
	})
}

func TestConsumerGroupHandler_ConcurrentMessages(t *testing.T) {
	t.Run("handles empty message stream", func(t *testing.T) {
		ready := make(chan bool)
		processor := &kafka.MessageValueProcessor{}
		handler := NewConsumerGroupHandler(ready, processor)

		session := NewMockConsumerGroupSession()
		claim := NewMockConsumerGroupClaim()

		go func() {
			// Just close without sending messages
			close(claim.messages)
		}()

		err := handler.ConsumeClaim(session, claim)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(session.markedMessages) != 0 {
			t.Errorf("Expected 0 marked messages, got %d", len(session.markedMessages))
		}
	})
}
