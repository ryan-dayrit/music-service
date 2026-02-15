package producer

import (
	"testing"

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
