package producer

import (
	"context"
	"testing"

	"music-service/gen/pb"
	"music-service/pkg/kafka"
)

func TestNewProducerHandler_InvalidBrokers(t *testing.T) {
	// Note: confluent-kafka-go creates producers even with empty/invalid brokers
	// Actual validation happens when trying to produce messages
	cfg := kafka.Config{
		Brokers: "",
		Topics:  "test-topic",
	}

	handler, err := NewProducerHandler(cfg)

	// Empty brokers creates producer but doesn't fail construction
	if err != nil {
		t.Errorf("NewProducerHandler() unexpected error = %v", err)
	}

	if handler == nil {
		t.Error("NewProducerHandler() should return non-nil handler")
	}
}

func TestNewProducerHandler_EmptyTopics(t *testing.T) {
	cfg := kafka.Config{
		Brokers: "localhost:9092",
		Topics:  "",
	}

	handler, err := NewProducerHandler(cfg)

	if err == nil {
		t.Error("NewProducerHandler() expected error for empty topics, got nil")
	}
	if handler != nil {
		t.Error("NewProducerHandler() should return nil handler on error")
	}
}

func TestNewProducerHandler_WhitespaceOnlyTopics(t *testing.T) {
	cfg := kafka.Config{
		Brokers: "localhost:9092",
		Topics:  "   ",
	}

	handler, err := NewProducerHandler(cfg)

	if err == nil {
		t.Error("NewProducerHandler() expected error for whitespace-only topics, got nil")
	}
	if handler != nil {
		t.Error("NewProducerHandler() should return nil handler on error")
	}
}

func TestNewProducerHandler_TopicWithLeadingTrailingSpaces(t *testing.T) {
	cfg := kafka.Config{
		Brokers: "localhost:9092",
		Topics:  "  test-topic  , other-topic",
	}

	handler, err := NewProducerHandler(cfg)

	if err != nil {
		t.Errorf("NewProducerHandler() unexpected error = %v", err)
	}
	if handler == nil {
		t.Error("NewProducerHandler() should return non-nil handler")
	}

	ph, ok := handler.(*producerHandler)
	if !ok {
		t.Fatal("handler is not *producerHandler")
	}
	if ph.topic != "test-topic" {
		t.Errorf("expected trimmed topic 'test-topic', got '%s'", ph.topic)
	}
}

func TestProducerHandler_Struct(t *testing.T) {
	ph := &producerHandler{
		cfg: kafka.Config{
			Brokers: "localhost:9092",
			Topics:  "test-topic",
		},
		confluentProducer: nil,
	}

	if ph.cfg.Brokers != "localhost:9092" {
		t.Errorf("Expected brokers 'localhost:9092', got '%s'", ph.cfg.Brokers)
	}

	if ph.cfg.Topics != "test-topic" {
		t.Errorf("Expected topics 'test-topic', got '%s'", ph.cfg.Topics)
	}

	if ph.confluentProducer != nil {
		t.Error("Expected nil producer")
	}
}

func TestProducerHandler_InterfaceCompliance(t *testing.T) {
	var _ kafka.ProducerHandler = (*producerHandler)(nil)
}

func TestProducerHandler_Produce_CancelledContext(t *testing.T) {
	// With a nil confluentProducer, calling Produce() with a cancelled context
	// should return early without panicking (context check happens before producer call).
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ph := &producerHandler{
		cfg: kafka.Config{
			Brokers: "localhost:9092",
			Topics:  "test-topic",
		},
		topic:             "test-topic",
		confluentProducer: nil,
	}

	// Should not panic even though confluentProducer is nil,
	// because ctx is already cancelled and we return early.
	ph.Produce(ctx, &pb.Album{Id: 1, Title: "Test Album"})
}

func TestNewProducerHandler_ValidBrokers(t *testing.T) {
	cfg := kafka.Config{
		Brokers: "localhost:9092",
		Topics:  "test-topic",
	}

	handler, err := NewProducerHandler(cfg)

	if err != nil {
		t.Errorf("NewProducerHandler() unexpected error = %v", err)
	}

	if handler == nil {
		t.Error("NewProducerHandler() should return non-nil handler")
	}
}

func TestNewProducerHandler_MultipleBrokers(t *testing.T) {
	cfg := kafka.Config{
		Brokers: "localhost:9092,localhost:9093,localhost:9094",
		Topics:  "test-topic",
	}

	handler, err := NewProducerHandler(cfg)

	if err != nil {
		t.Errorf("NewProducerHandler() unexpected error = %v", err)
	}

	if handler == nil {
		t.Error("NewProducerHandler() should return non-nil handler")
	}
}

func TestProducerHandler_ConfigStorage(t *testing.T) {
	cfg := kafka.Config{
		Brokers: "kafka1:9092,kafka2:9092",
		Topics:  "music.albums",
	}

	ph := &producerHandler{
		cfg:               cfg,
		confluentProducer: nil,
	}

	if ph.cfg.Brokers != cfg.Brokers {
		t.Errorf("Brokers mismatch: got %s, want %s", ph.cfg.Brokers, cfg.Brokers)
	}

	if ph.cfg.Topics != cfg.Topics {
		t.Errorf("Topics mismatch: got %s, want %s", ph.cfg.Topics, cfg.Topics)
	}
}
