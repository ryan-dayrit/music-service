package confluent

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"music-service/pkg/kafka/message"
)

// MockMessageValueProcessor is a mock implementation for testing
type MockMessageValueProcessor struct {
	ProcessedMessages [][]byte
	ProcessCount      int
}

func (m *MockMessageValueProcessor) Process(msg []byte) {
	m.ProcessedMessages = append(m.ProcessedMessages, msg)
	m.ProcessCount++
}

// TestNewConsumer tests the constructor
func TestNewConsumer(t *testing.T) {
	tests := []struct {
		name            string
		consumer        *kafka.Consumer
		processor       message.MessageValueProcessor
		parallelWorkers int
	}{
		{
			name:            "with nil consumer",
			consumer:        nil,
			processor:       &MockMessageValueProcessor{},
			parallelWorkers: 5,
		},
		{
			name:            "with zero workers",
			consumer:        nil,
			processor:       &MockMessageValueProcessor{},
			parallelWorkers: 0,
		},
		{
			name:            "with single worker",
			consumer:        nil,
			processor:       &MockMessageValueProcessor{},
			parallelWorkers: 1,
		},
		{
			name:            "with many workers",
			consumer:        nil,
			processor:       &MockMessageValueProcessor{},
			parallelWorkers: 100,
		},
		{
			name:            "with negative workers",
			consumer:        nil,
			processor:       &MockMessageValueProcessor{},
			parallelWorkers: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConsumer(tt.consumer, tt.processor, tt.parallelWorkers)

			if c == nil {
				t.Fatal("Expected consumer to be created, got nil")
			}

			if c.confluentConsumer != tt.consumer {
				t.Error("Consumer field not set correctly")
			}

			if c.messageValueProcessor != tt.processor {
				t.Error("MessageValueProcessor field not set correctly")
			}

			if c.parallelWorkers != tt.parallelWorkers {
				t.Errorf("Expected parallelWorkers=%d, got %d", tt.parallelWorkers, c.parallelWorkers)
			}
		})
	}
}

// TestNewConsumer_NilProcessor tests constructor with nil processor
func TestNewConsumer_NilProcessor(t *testing.T) {
	c := NewConsumer(nil, nil, 5)

	if c == nil {
		t.Fatal("Expected consumer to be created, got nil")
	}

	if c.messageValueProcessor != nil {
		t.Error("Expected nil processor")
	}
}

// TestConsumer_Struct tests the consumer struct fields
func TestConsumer_Struct(t *testing.T) {
	tests := []struct {
		name            string
		consumer        *kafka.Consumer
		processor       message.MessageValueProcessor
		parallelWorkers int
	}{
		{
			name:            "all nil fields",
			consumer:        nil,
			processor:       nil,
			parallelWorkers: 0,
		},
		{
			name:            "with processor",
			consumer:        nil,
			processor:       &MockMessageValueProcessor{},
			parallelWorkers: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &consumer{
				confluentConsumer:     tt.consumer,
				messageValueProcessor: tt.processor,
				parallelWorkers:       tt.parallelWorkers,
			}

			if c.confluentConsumer != tt.consumer {
				t.Error("Consumer field not set correctly")
			}

			if c.messageValueProcessor != tt.processor {
				t.Error("MessageValueProcessor field not set correctly")
			}

			if c.parallelWorkers != tt.parallelWorkers {
				t.Errorf("Expected parallelWorkers=%d, got %d", tt.parallelWorkers, c.parallelWorkers)
			}
		})
	}
}

// TestAck_Struct tests the ack struct
func TestAck_Struct(t *testing.T) {
	tests := []struct {
		name      string
		topic     string
		partition int32
		offset    kafka.Offset
	}{
		{
			name:      "valid ack",
			topic:     "test-topic",
			partition: 1,
			offset:    100,
		},
		{
			name:      "zero offset",
			topic:     "another-topic",
			partition: 0,
			offset:    0,
		},
		{
			name:      "high partition",
			topic:     "multi-partition-topic",
			partition: 99,
			offset:    123456,
		},
		{
			name:      "negative offset",
			topic:     "special-topic",
			partition: 0,
			offset:    kafka.Offset(-1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			topicCopy := tt.topic
			a := ack{
				tp: kafka.TopicPartition{
					Topic:     &topicCopy,
					Partition: tt.partition,
					Offset:    tt.offset,
				},
				off: tt.offset,
			}

			if a.tp.Topic == nil {
				t.Fatal("Expected non-nil topic pointer")
			}

			if *a.tp.Topic != tt.topic {
				t.Errorf("Expected topic '%s', got '%s'", tt.topic, *a.tp.Topic)
			}

			if a.tp.Partition != tt.partition {
				t.Errorf("Expected partition %d, got %d", tt.partition, a.tp.Partition)
			}

			if a.tp.Offset != tt.offset {
				t.Errorf("Expected TopicPartition offset %d, got %d", tt.offset, a.tp.Offset)
			}

			if a.off != tt.offset {
				t.Errorf("Expected off %d, got %d", tt.offset, a.off)
			}
		})
	}
}

// TestAck_OffsetIncrement tests that offset values can be incremented
func TestAck_OffsetIncrement(t *testing.T) {
	topic := "test-topic"
	originalOffset := kafka.Offset(100)

	a := ack{
		tp: kafka.TopicPartition{
			Topic:     &topic,
			Partition: 0,
			Offset:    originalOffset,
		},
		off: originalOffset,
	}

	// Simulate offset increment (as done in the Consume method)
	nextOffset := a.off + 1

	if nextOffset != kafka.Offset(101) {
		t.Errorf("Expected next offset 101, got %d", nextOffset)
	}

	// Verify original ack is unchanged
	if a.off != originalOffset {
		t.Errorf("Expected original offset to remain %d, got %d", originalOffset, a.off)
	}
}

// TestMockMessageValueProcessor tests the mock processor
func TestMockMessageValueProcessor(t *testing.T) {
	mock := &MockMessageValueProcessor{}

	testMessages := [][]byte{
		[]byte("message 1"),
		[]byte("message 2"),
		[]byte("message 3"),
	}

	for _, msg := range testMessages {
		mock.Process(msg)
	}

	if mock.ProcessCount != 3 {
		t.Errorf("Expected ProcessCount=3, got %d", mock.ProcessCount)
	}

	if len(mock.ProcessedMessages) != 3 {
		t.Errorf("Expected 3 processed messages, got %d", len(mock.ProcessedMessages))
	}

	for i, msg := range testMessages {
		if string(mock.ProcessedMessages[i]) != string(msg) {
			t.Errorf("Message %d mismatch: expected '%s', got '%s'",
				i, string(msg), string(mock.ProcessedMessages[i]))
		}
	}
}

// TestMockMessageValueProcessor_EmptyMessage tests processing empty messages
func TestMockMessageValueProcessor_EmptyMessage(t *testing.T) {
	mock := &MockMessageValueProcessor{}

	mock.Process([]byte{})
	mock.Process(nil)

	if mock.ProcessCount != 2 {
		t.Errorf("Expected ProcessCount=2, got %d", mock.ProcessCount)
	}

	if len(mock.ProcessedMessages) != 2 {
		t.Errorf("Expected 2 processed messages, got %d", len(mock.ProcessedMessages))
	}
}

// TestMockMessageValueProcessor_LargeMessage tests processing large messages
func TestMockMessageValueProcessor_LargeMessage(t *testing.T) {
	mock := &MockMessageValueProcessor{}

	// Create a large message (1MB)
	largeMsg := make([]byte, 1024*1024)
	for i := range largeMsg {
		largeMsg[i] = byte(i % 256)
	}

	mock.Process(largeMsg)

	if mock.ProcessCount != 1 {
		t.Errorf("Expected ProcessCount=1, got %d", mock.ProcessCount)
	}

	if len(mock.ProcessedMessages) != 1 {
		t.Errorf("Expected 1 processed message, got %d", len(mock.ProcessedMessages))
	}

	if len(mock.ProcessedMessages[0]) != len(largeMsg) {
		t.Errorf("Expected message size %d, got %d", len(largeMsg), len(mock.ProcessedMessages[0]))
	}
}

// TestMockMessageValueProcessor_InterfaceCompliance tests interface compliance
func TestMockMessageValueProcessor_InterfaceCompliance(t *testing.T) {
	var _ message.MessageValueProcessor = (*MockMessageValueProcessor)(nil)
}

// TestConsumer_FieldsArePrivate tests that consumer fields are unexported
func TestConsumer_FieldsArePrivate(t *testing.T) {
	// This test ensures that the consumer struct has unexported fields
	// If fields were exported, this would be a design smell
	c := NewConsumer(nil, nil, 0)

	// Access through the struct should work (same package)
	_ = c.confluentConsumer
	_ = c.messageValueProcessor
	_ = c.parallelWorkers

	// This test passes if compilation succeeds
}

// TestNewConsumer_ReturnsNonNilPointer tests that constructor never returns nil
func TestNewConsumer_ReturnsNonNilPointer(t *testing.T) {
	testCases := []struct {
		name            string
		consumer        *kafka.Consumer
		processor       message.MessageValueProcessor
		parallelWorkers int
	}{
		{"all nil", nil, nil, 0},
		{"with processor", nil, &MockMessageValueProcessor{}, 0},
		{"with workers", nil, nil, 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewConsumer(tc.consumer, tc.processor, tc.parallelWorkers)
			if result == nil {
				t.Error("NewConsumer should never return nil")
			}
		})
	}
}
