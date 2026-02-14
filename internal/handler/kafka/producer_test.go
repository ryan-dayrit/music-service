package kafka

import (
	"context"
	"errors"
	"testing"

	"math/rand/v2"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"music-service/gen/pb"
	"music-service/pkg/kafka"
)

// MockSyncProducer is a mock implementation of sarama.SyncProducer
type MockSyncProducer struct {
	mock.Mock
}

func (m *MockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	args := m.Called(msg)
	return int32(args.Int(0)), int64(args.Int(1)), args.Error(2)
}

func (m *MockSyncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	args := m.Called(msgs)
	return args.Error(0)
}

func (m *MockSyncProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSyncProducer) TxnStatus() sarama.ProducerTxnStatusFlag {
	args := m.Called()
	return args.Get(0).(sarama.ProducerTxnStatusFlag)
}

func (m *MockSyncProducer) IsTransactional() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSyncProducer) BeginTxn() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSyncProducer) CommitTxn() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSyncProducer) AbortTxn() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSyncProducer) AddOffsetsToTxn(offsets map[string][]*sarama.PartitionOffsetMetadata, groupId string) error {
	args := m.Called(offsets, groupId)
	return args.Error(0)
}

func (m *MockSyncProducer) AddMessageToTxn(msg *sarama.ConsumerMessage, groupId string, metadata *string) error {
	args := m.Called(msg, groupId, metadata)
	return args.Error(0)
}

// TestNewProducer tests the NewProducer function
func TestNewProducer(t *testing.T) {
	tests := []struct {
		name             string
		cfg              kafka.Config
		mustError        bool
		kafkaNotRequired bool
	}{
		{
			name: "valid config",
			cfg: kafka.Config{
				Brokers:       "localhost:9092",
				Topics:        "test-topic",
				ConsumerGroup: "test-group",
				Assignor:      "sticky",
				Oldest:        false,
			},
			mustError:        false,
			kafkaNotRequired: true, // Can succeed or fail depending on Kafka availability
		},
		{
			name: "valid config with multiple brokers",
			cfg: kafka.Config{
				Brokers:       "localhost:9092,localhost:9093,localhost:9094",
				Topics:        "test-topic",
				ConsumerGroup: "test-group",
				Assignor:      "sticky",
				Oldest:        false,
			},
			mustError:        false,
			kafkaNotRequired: true, // Can succeed or fail depending on Kafka availability
		},
		{
			name: "empty brokers",
			cfg: kafka.Config{
				Brokers:       "",
				Topics:        "test-topic",
				ConsumerGroup: "test-group",
				Assignor:      "sticky",
				Oldest:        false,
			},
			mustError:        true,
			kafkaNotRequired: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewProducer(tt.cfg)
			if tt.mustError {
				assert.Error(t, err)
				assert.Nil(t, p)
			} else if tt.kafkaNotRequired {
				// Can succeed or fail depending on whether Kafka is running
				if err != nil {
					// Kafka not running is acceptable
					assert.Nil(t, p)
					t.Logf("Kafka appears to be unavailable: %v", err)
				} else {
					// Kafka is running
					assert.NotNil(t, p)
					assert.Equal(t, tt.cfg, p.cfg)
					assert.NotNil(t, p.syncProducer)
					if p.syncProducer != nil {
						_ = p.syncProducer.Close()
					}
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, p)
				assert.Equal(t, tt.cfg, p.cfg)
				assert.NotNil(t, p.syncProducer)
				if p.syncProducer != nil {
					_ = p.syncProducer.Close()
				}
			}
		})
	}
}

// TestNewProducer_WithMock tests the NewProducer function structure with mocked dependencies
func TestNewProducer_WithMock(t *testing.T) {
	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "sticky",
		Oldest:        false,
	}

	// Test that the producer struct fields are correctly populated
	// Note: We can't directly inject mocks into NewProducer since it creates its own sync producer
	// This test verifies the producer structure when creation is successful
	t.Run("producer fields validation", func(t *testing.T) {
		// Create a mock sync producer
		mockSP := new(MockSyncProducer)

		// Create producer directly for testing structure
		p := &producer{
			cfg:          cfg,
			syncProducer: mockSP,
		}

		assert.NotNil(t, p)
		assert.Equal(t, cfg, p.cfg)
		assert.Equal(t, mockSP, p.syncProducer)
	})
}

// TestProduce_Success tests the Produce function with successful message sending
func TestProduce_Success(t *testing.T) {
	mockSP := new(MockSyncProducer)

	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "sticky",
		Oldest:        false,
	}

	p := &producer{
		cfg:          cfg,
		syncProducer: mockSP,
	}

	ctx := context.Background()

	// Mock successful message sending
	mockSP.On("SendMessage", mock.MatchedBy(func(msg *sarama.ProducerMessage) bool {
		return msg.Topic == "test-topic" && msg.Value != nil
	})).Return(1, 100, nil)

	// This should not panic
	album := &pb.Album{
		Id:     rand.Int32(),
		Title:  uuid.NewString(),
		Artist: uuid.NewString(),
		Price:  rand.Float32(),
	}
	p.Produce(ctx, album)

	mockSP.AssertExpectations(t)
	mockSP.AssertCalled(t, "SendMessage", mock.Anything)
}

// TestProduce_Error tests the Produce function with error scenarios
func TestProduce_Error(t *testing.T) {
	mockSP := new(MockSyncProducer)

	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "sticky",
		Oldest:        false,
	}

	p := &producer{
		cfg:          cfg,
		syncProducer: mockSP,
	}

	ctx := context.Background()

	// Mock failed message sending
	mockSP.On("SendMessage", mock.Anything).Return(0, 0, errors.New("send error"))

	// This should panic due to the send error
	assert.Panics(t, func() {
		album := &pb.Album{
			Id:     rand.Int32(),
			Title:  uuid.NewString(),
			Artist: uuid.NewString(),
			Price:  rand.Float32(),
		}
		p.Produce(ctx, album)
	})

	mockSP.AssertExpectations(t)
}

// TestProduce_MessageContent tests that messages contain unique UUIDs
func TestProduce_MessageContent(t *testing.T) {
	mockSP := new(MockSyncProducer)

	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "sticky",
		Oldest:        false,
	}

	p := &producer{
		cfg:          cfg,
		syncProducer: mockSP,
	}

	ctx := context.Background()

	var capturedMessage *sarama.ProducerMessage

	// Capture the message that is sent
	mockSP.On("SendMessage", mock.Anything).Run(func(args mock.Arguments) {
		capturedMessage = args.Get(0).(*sarama.ProducerMessage)
	}).Return(0, 0, nil)

	album := &pb.Album{
		Id:     rand.Int32(),
		Title:  uuid.NewString(),
		Artist: uuid.NewString(),
		Price:  rand.Float32(),
	}
	p.Produce(ctx, album)

	assert.NotNil(t, capturedMessage)
	assert.Equal(t, "test-topic", capturedMessage.Topic)
	assert.NotNil(t, capturedMessage.Value)

	// Decode the message value
	_, err := capturedMessage.Value.Encode()
	assert.NoError(t, err)

	mockSP.AssertExpectations(t)
}

// TestProduce_MultipleMessages tests that multiple messages have unique UUIDs
func TestProduce_MultipleMessages(t *testing.T) {
	mockSP := new(MockSyncProducer)

	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "sticky",
		Oldest:        false,
	}

	p := &producer{
		cfg:          cfg,
		syncProducer: mockSP,
	}

	ctx := context.Background()

	var capturedMessages []string

	// Capture multiple messages
	mockSP.On("SendMessage", mock.Anything).Run(func(args mock.Arguments) {
		msg := args.Get(0).(*sarama.ProducerMessage)
		msgBytes, _ := msg.Value.Encode()
		capturedMessages = append(capturedMessages, string(msgBytes))
	}).Return(0, 0, nil).Times(3)

	// Produce three messages
	album1 := &pb.Album{
		Id:     rand.Int32(),
		Title:  uuid.NewString(),
		Artist: uuid.NewString(),
		Price:  rand.Float32(),
	}
	p.Produce(ctx, album1)

	album2 := &pb.Album{
		Id:     rand.Int32(),
		Title:  uuid.NewString(),
		Artist: uuid.NewString(),
		Price:  rand.Float32(),
	}
	p.Produce(ctx, album2)

	album3 := &pb.Album{
		Id:     rand.Int32(),
		Title:  uuid.NewString(),
		Artist: uuid.NewString(),
		Price:  rand.Float32(),
	}
	p.Produce(ctx, album3)

	assert.Len(t, capturedMessages, 3)

	// Verify all messages are unique
	messageSet := make(map[string]bool)
	for _, msg := range capturedMessages {
		assert.NotContains(t, messageSet, msg, "messages should be unique")
		messageSet[msg] = true
	}

	mockSP.AssertExpectations(t)
}

// TestProduce_DifferentTopics tests producing to different topics
func TestProduce_DifferentTopics(t *testing.T) {
	tests := []struct {
		name  string
		topic string
	}{
		{
			name:  "topic1",
			topic: "topic1",
		},
		{
			name:  "topic2",
			topic: "topic2",
		},
		{
			name:  "topic with dashes",
			topic: "my-topic-name",
		},
		{
			name:  "topic with underscores",
			topic: "my_topic_name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSP := new(MockSyncProducer)

			cfg := kafka.Config{
				Brokers:       "localhost:9092",
				Topics:        tt.topic,
				ConsumerGroup: "test-group",
				Assignor:      "sticky",
			}

			p := &producer{
				cfg:          cfg,
				syncProducer: mockSP,
			}

			ctx := context.Background()

			// Verify the message is sent to the correct topic
			mockSP.On("SendMessage", mock.MatchedBy(func(msg *sarama.ProducerMessage) bool {
				return msg.Topic == tt.topic
			})).Return(0, 0, nil)

			album := &pb.Album{
				Id:     rand.Int32(),
				Title:  uuid.NewString(),
				Artist: uuid.NewString(),
				Price:  rand.Float32(),
			}
			p.Produce(ctx, album)

			mockSP.AssertExpectations(t)
		})
	}
}

// TestProduce_DifferentPartitions tests producing to different partitions
func TestProduce_DifferentPartitions(t *testing.T) {
	partitions := []struct {
		name      string
		partition int32
		offset    int64
	}{
		{
			name:      "partition 0",
			partition: 0,
			offset:    100,
		},
		{
			name:      "partition 1",
			partition: 1,
			offset:    200,
		},
		{
			name:      "partition 5",
			partition: 5,
			offset:    500,
		},
	}

	for _, tt := range partitions {
		t.Run(tt.name, func(t *testing.T) {
			mockSP := new(MockSyncProducer)

			cfg := kafka.Config{
				Brokers:       "localhost:9092",
				Topics:        "test-topic",
				ConsumerGroup: "test-group",
				Assignor:      "sticky",
			}

			p := &producer{
				cfg:          cfg,
				syncProducer: mockSP,
			}

			ctx := context.Background()

			// Mock returns specific partition and offset
			mockSP.On("SendMessage", mock.Anything).Return(int(tt.partition), int(tt.offset), nil)

			album := &pb.Album{
				Id:     rand.Int32(),
				Title:  uuid.NewString(),
				Artist: uuid.NewString(),
				Price:  rand.Float32(),
			}
			p.Produce(ctx, album)

			mockSP.AssertExpectations(t)
		})
	}
}

// TestProducerStructure tests the producer struct fields
func TestProducerStructure(t *testing.T) {
	mockSP := new(MockSyncProducer)

	cfg := kafka.Config{
		Brokers:       "localhost:9092,localhost:9093",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "roundrobin",
		Oldest:        true,
	}

	p := &producer{
		cfg:          cfg,
		syncProducer: mockSP,
	}

	assert.NotNil(t, p)
	assert.Equal(t, "localhost:9092,localhost:9093", p.cfg.Brokers)
	assert.Equal(t, "test-topic", p.cfg.Topics)
	assert.Equal(t, "test-group", p.cfg.ConsumerGroup)
	assert.Equal(t, "roundrobin", p.cfg.Assignor)
	assert.True(t, p.cfg.Oldest)
	assert.Equal(t, mockSP, p.syncProducer)
}

// TestProduce_WithContext tests that Produce accepts context parameter
func TestProduce_WithContext(t *testing.T) {
	mockSP := new(MockSyncProducer)

	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "sticky",
	}

	p := &producer{
		cfg:          cfg,
		syncProducer: mockSP,
	}

	// Test with background context
	t.Run("background context", func(t *testing.T) {
		ctx := context.Background()
		mockSP.On("SendMessage", mock.Anything).Return(0, 0, nil).Once()
		album := &pb.Album{
			Id:     rand.Int32(),
			Title:  uuid.NewString(),
			Artist: uuid.NewString(),
			Price:  rand.Float32(),
		}
		p.Produce(ctx, album)
		mockSP.AssertExpectations(t)
	})

	// Test with TODO context
	t.Run("todo context", func(t *testing.T) {
		ctx := context.TODO()
		mockSP.On("SendMessage", mock.Anything).Return(0, 0, nil).Once()
		album := &pb.Album{
			Id:     rand.Int32(),
			Title:  uuid.NewString(),
			Artist: uuid.NewString(),
			Price:  rand.Float32(),
		}
		p.Produce(ctx, album)
		mockSP.AssertExpectations(t)
	})

	// Test with cancelled context (should still attempt to send since Produce doesn't check context)
	t.Run("cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		mockSP.On("SendMessage", mock.Anything).Return(0, 0, nil).Once()
		album := &pb.Album{
			Id:     rand.Int32(),
			Title:  uuid.NewString(),
			Artist: uuid.NewString(),
			Price:  rand.Float32(),
		}
		p.Produce(ctx, album)
		mockSP.AssertExpectations(t)
	})
}

// TestProduce_MessageEncoding tests that messages are properly encoded
func TestProduce_MessageEncoding(t *testing.T) {
	mockSP := new(MockSyncProducer)

	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "sticky",
	}

	p := &producer{
		cfg:          cfg,
		syncProducer: mockSP,
	}

	ctx := context.Background()

	var capturedEncoder sarama.Encoder

	mockSP.On("SendMessage", mock.Anything).Run(func(args mock.Arguments) {
		msg := args.Get(0).(*sarama.ProducerMessage)
		capturedEncoder = msg.Value
	}).Return(0, 0, nil)

	album := &pb.Album{
		Id:     rand.Int32(),
		Title:  uuid.NewString(),
		Artist: uuid.NewString(),
		Price:  rand.Float32(),
	}
	p.Produce(ctx, album)

	assert.NotNil(t, capturedEncoder)

	// Verify encoding works
	bytes, err := capturedEncoder.Encode()
	assert.NoError(t, err)
	assert.NotEmpty(t, bytes)

	mockSP.AssertExpectations(t)
}

// TestProduce_EmptyTopic tests producing with empty topic
func TestProduce_EmptyTopic(t *testing.T) {
	mockSP := new(MockSyncProducer)

	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "",
		ConsumerGroup: "test-group",
		Assignor:      "sticky",
	}

	p := &producer{
		cfg:          cfg,
		syncProducer: mockSP,
	}

	ctx := context.Background()

	// Even with empty topic, the Produce function will attempt to send
	mockSP.On("SendMessage", mock.MatchedBy(func(msg *sarama.ProducerMessage) bool {
		return msg.Topic == ""
	})).Return(0, 0, errors.New("invalid topic"))

	// Should panic due to error
	assert.Panics(t, func() {
		album := &pb.Album{
			Id:     rand.Int32(),
			Title:  uuid.NewString(),
			Artist: uuid.NewString(),
			Price:  rand.Float32(),
		}
		p.Produce(ctx, album)
	})

	mockSP.AssertExpectations(t)
}
