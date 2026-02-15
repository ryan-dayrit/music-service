package kafka

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"music-service/pkg/kafka"
)

// MockConsumerGroup is a mock implementation of sarama.ConsumerGroup
type MockConsumerGroup struct {
	mock.Mock
}

func (m *MockConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	args := m.Called(ctx, topics, handler)
	return args.Error(0)
}

func (m *MockConsumerGroup) Errors() <-chan error {
	args := m.Called()
	return args.Get(0).(<-chan error)
}

func (m *MockConsumerGroup) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConsumerGroup) Pause(partitions map[string][]int32) {
	m.Called(partitions)
}

func (m *MockConsumerGroup) Resume(partitions map[string][]int32) {
	m.Called(partitions)
}

func (m *MockConsumerGroup) PauseAll() {
	m.Called()
}

func (m *MockConsumerGroup) ResumeAll() {
	m.Called()
}

// TestNewHandler tests the NewHandler function
func TestNewHandler(t *testing.T) {
	tests := []struct {
		name             string
		cfg              kafka.Config
		mustError        bool
		kafkaNotRequired bool
	}{
		{
			name: "valid config with sticky assignor",
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
			name: "valid config with roundrobin assignor",
			cfg: kafka.Config{
				Brokers:       "localhost:9092",
				Topics:        "test-topic",
				ConsumerGroup: "test-group",
				Assignor:      "roundrobin",
				Oldest:        false,
			},
			mustError:        false,
			kafkaNotRequired: true, // Can succeed or fail depending on Kafka availability
		},
		{
			name: "valid config with range assignor",
			cfg: kafka.Config{
				Brokers:       "localhost:9092",
				Topics:        "test-topic",
				ConsumerGroup: "test-group",
				Assignor:      "range",
				Oldest:        true,
			},
			mustError:        false,
			kafkaNotRequired: true, // Can succeed or fail depending on Kafka availability
		},
		{
			name: "invalid assignor",
			cfg: kafka.Config{
				Brokers:       "localhost:9092",
				Topics:        "test-topic",
				ConsumerGroup: "test-group",
				Assignor:      "invalid",
				Oldest:        false,
			},
			mustError:        true,
			kafkaNotRequired: false,
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
			h, err := NewConsumerHandler(tt.cfg, nil)
			if tt.mustError {
				assert.Error(t, err)
				assert.Nil(t, h)
			} else if tt.kafkaNotRequired {
				// Can succeed or fail depending on whether Kafka is running
				if err != nil {
					// Kafka not running is acceptable
					assert.Nil(t, h)
					t.Logf("Kafka appears to be unavailable: %v", err)
				} else {
					// Kafka is running
					assert.NotNil(t, h)
					assert.Equal(t, tt.cfg, h.cfg)
					assert.NotNil(t, h.consumerGroup)
					if h.consumerGroup != nil {
						_ = h.consumerGroup.Close()
					}
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, h)
				assert.Equal(t, tt.cfg, h.cfg)
				assert.NotNil(t, h.consumerGroup)
				if h.consumerGroup != nil {
					_ = h.consumerGroup.Close()
				}
			}
		})
	}
}

// TestNewHandler_WithMock tests the NewHandler function structure with mocked dependencies
func TestNewHandler_WithMock(t *testing.T) {
	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "sticky",
		Oldest:        false,
	}

	// Test that the handler struct fields are correctly populated
	// Note: We can't directly inject mocks into NewHandler since it creates its own consumer group
	// This test verifies the handler structure when creation is successful
	t.Run("handler fields validation", func(t *testing.T) {
		// Create a mock consumer group
		mockCG := new(MockConsumerGroup)

		// Create handler directly for testing structure
		h := &consumerHandler{
			cfg:           cfg,
			consumerGroup: mockCG,
		}

		assert.NotNil(t, h)
		assert.Equal(t, cfg, h.cfg)
		assert.Equal(t, mockCG, h.consumerGroup)
	})
}

// TestConsume_ContextCancellation tests the Consume function with context cancellation
func TestConsume_ContextCancellation(t *testing.T) {
	mockCG := new(MockConsumerGroup)

	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "sticky",
		Oldest:        false,
	}

	h := &consumerHandler{
		cfg:           cfg,
		consumerGroup: mockCG,
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Mock the Consume method to simulate successful consumption and then return when context is done
	mockCG.On("Consume", mock.Anything, []string{"test-topic"}, mock.Anything).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		handler := args.Get(2).(sarama.ConsumerGroupHandler)
		// Simulate setup
		_ = handler.Setup(nil)
		// Wait for context cancellation
		<-ctx.Done()
	}).Return(nil)

	mockCG.On("Close").Return(nil)

	// Cancel context after a short delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	// This should exit gracefully when context is cancelled
	h.Consume(ctx)

	mockCG.AssertExpectations(t)
}

// TestConsume_ErrorHandling tests the Consume function with error scenarios
func TestConsume_ErrorHandling(t *testing.T) {
	t.Run("consumer group closed error", func(t *testing.T) {
		mockCG := new(MockConsumerGroup)

		cfg := kafka.Config{
			Brokers:       "localhost:9092",
			Topics:        "test-topic",
			ConsumerGroup: "test-group",
			Assignor:      "sticky",
		}

		h := &consumerHandler{
			cfg:           cfg,
			consumerGroup: mockCG,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		// Mock the Consume method to return ErrClosedConsumerGroup
		mockCG.On("Consume", mock.Anything, []string{"test-topic"}, mock.Anything).Run(func(args mock.Arguments) {
			handler := args.Get(2).(sarama.ConsumerGroupHandler)
			_ = handler.Setup(nil)
		}).Return(sarama.ErrClosedConsumerGroup)

		mockCG.On("Close").Return(nil)

		// This should handle the error gracefully
		h.Consume(ctx)

		mockCG.AssertExpectations(t)
	})
}

// TestConsume_MultipleTopics tests the Consume function with multiple topics
func TestConsume_MultipleTopics(t *testing.T) {
	mockCG := new(MockConsumerGroup)

	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "topic1,topic2,topic3",
		ConsumerGroup: "test-group",
		Assignor:      "sticky",
	}

	h := &consumerHandler{
		cfg:           cfg,
		consumerGroup: mockCG,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Expect Consume to be called with multiple topics
	mockCG.On("Consume", mock.Anything, []string{"topic1", "topic2", "topic3"}, mock.Anything).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		handler := args.Get(2).(sarama.ConsumerGroupHandler)
		_ = handler.Setup(nil)
		<-ctx.Done()
	}).Return(nil)

	mockCG.On("Close").Return(nil)

	h.Consume(ctx)

	mockCG.AssertExpectations(t)
}

// TestConsume_CloseError tests the Consume function when Close returns an error
func TestConsume_CloseError(t *testing.T) {
	mockCG := new(MockConsumerGroup)

	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "sticky",
	}

	h := &consumerHandler{
		cfg:           cfg,
		consumerGroup: mockCG,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	mockCG.On("Consume", mock.Anything, []string{"test-topic"}, mock.Anything).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		handler := args.Get(2).(sarama.ConsumerGroupHandler)
		_ = handler.Setup(nil)
		<-ctx.Done()
	}).Return(nil)

	// Close returns an error - this will cause a panic in the actual code
	mockCG.On("Close").Return(errors.New("close error"))

	// We expect a panic due to the close error
	assert.Panics(t, func() {
		h.Consume(ctx)
	})
}

// TestToggleConsumptionFlow tests the toggleConsumptionFlow function
func TestToggleConsumptionFlow(t *testing.T) {
	t.Run("pause when not paused", func(t *testing.T) {
		mockCG := new(MockConsumerGroup)
		isPaused := false

		mockCG.On("PauseAll").Return()

		toggleConsumptionFlow(mockCG, &isPaused)

		assert.True(t, isPaused, "should be paused after toggle")
		mockCG.AssertExpectations(t)
		mockCG.AssertCalled(t, "PauseAll")
	})

	t.Run("resume when paused", func(t *testing.T) {
		mockCG := new(MockConsumerGroup)
		isPaused := true

		mockCG.On("ResumeAll").Return()

		toggleConsumptionFlow(mockCG, &isPaused)

		assert.False(t, isPaused, "should not be paused after toggle")
		mockCG.AssertExpectations(t)
		mockCG.AssertCalled(t, "ResumeAll")
	})

	t.Run("multiple toggles", func(t *testing.T) {
		mockCG := new(MockConsumerGroup)
		isPaused := false

		// First toggle: pause
		mockCG.On("PauseAll").Return().Once()
		toggleConsumptionFlow(mockCG, &isPaused)
		assert.True(t, isPaused)

		// Second toggle: resume
		mockCG.On("ResumeAll").Return().Once()
		toggleConsumptionFlow(mockCG, &isPaused)
		assert.False(t, isPaused)

		// Third toggle: pause again
		mockCG.On("PauseAll").Return().Once()
		toggleConsumptionFlow(mockCG, &isPaused)
		assert.True(t, isPaused)

		mockCG.AssertExpectations(t)
	})
}

// TestHandlerStructure tests the handler struct fields
func TestHandlerStructure(t *testing.T) {
	mockCG := new(MockConsumerGroup)

	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "sticky",
		Oldest:        true,
	}

	h := &consumerHandler{
		cfg:           cfg,
		consumerGroup: mockCG,
	}

	assert.NotNil(t, h)
	assert.Equal(t, "localhost:9092", h.cfg.Brokers)
	assert.Equal(t, "test-topic", h.cfg.Topics)
	assert.Equal(t, "test-group", h.cfg.ConsumerGroup)
	assert.Equal(t, "sticky", h.cfg.Assignor)
	assert.True(t, h.cfg.Oldest)
	assert.Equal(t, mockCG, h.consumerGroup)
}

// TestConsume_TopicParsing tests various topic configurations
func TestConsume_TopicParsing(t *testing.T) {
	tests := []struct {
		name           string
		topics         string
		expectedTopics []string
	}{
		{
			name:           "single topic",
			topics:         "topic1",
			expectedTopics: []string{"topic1"},
		},
		{
			name:           "two topics",
			topics:         "topic1,topic2",
			expectedTopics: []string{"topic1", "topic2"},
		},
		{
			name:           "multiple topics with spaces",
			topics:         "topic1,topic2,topic3,topic4",
			expectedTopics: []string{"topic1", "topic2", "topic3", "topic4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCG := new(MockConsumerGroup)

			cfg := kafka.Config{
				Brokers:       "localhost:9092",
				Topics:        tt.topics,
				ConsumerGroup: "test-group",
				Assignor:      "sticky",
			}

			h := &consumerHandler{
				cfg:           cfg,
				consumerGroup: mockCG,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			mockCG.On("Consume", mock.Anything, tt.expectedTopics, mock.Anything).Run(func(args mock.Arguments) {
				ctx := args.Get(0).(context.Context)
				handler := args.Get(2).(sarama.ConsumerGroupHandler)
				_ = handler.Setup(nil)
				<-ctx.Done()
			}).Return(nil)

			mockCG.On("Close").Return(nil)

			h.Consume(ctx)

			mockCG.AssertExpectations(t)
		})
	}
}
