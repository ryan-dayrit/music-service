package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"

	"music-service/gen/pb"
)

type MockConsumerGroupSession struct {
	mock.Mock
}

func (m *MockConsumerGroupSession) Claims() map[string][]int32 {
	args := m.Called()
	return args.Get(0).(map[string][]int32)
}

func (m *MockConsumerGroupSession) MemberID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConsumerGroupSession) GenerationID() int32 {
	args := m.Called()
	return args.Get(0).(int32)
}

func (m *MockConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
	m.Called(topic, partition, offset, metadata)
}

func (m *MockConsumerGroupSession) Commit() {
	m.Called()
}

func (m *MockConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
	m.Called(topic, partition, offset, metadata)
}

func (m *MockConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	m.Called(msg, metadata)
}

func (m *MockConsumerGroupSession) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

type MockConsumerGroupClaim struct {
	mock.Mock
	messages chan *sarama.ConsumerMessage
}

func (m *MockConsumerGroupClaim) Topic() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConsumerGroupClaim) Partition() int32 {
	args := m.Called()
	return args.Get(0).(int32)
}

func (m *MockConsumerGroupClaim) InitialOffset() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockConsumerGroupClaim) HighWaterMarkOffset() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	return m.messages
}

func TestConsumerGroupHandler_Setup(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful setup",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cgh := &consumerGroupHandler{
				Ready: make(chan bool),
			}

			mockSession := new(MockConsumerGroupSession)

			go func() {
				err := cgh.Setup(mockSession)
				assert.NoError(t, err)
			}()

			select {
			case <-cgh.Ready:
			case <-time.After(1 * time.Second):
				t.Fatal("Ready channel was not closed within timeout")
			}
		})
	}
}

func TestConsumerGroupHandler_Cleanup(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful cleanup",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cgh := &consumerGroupHandler{
				Ready: make(chan bool),
			}

			mockSession := new(MockConsumerGroupSession)

			err := cgh.Cleanup(mockSession)
			assert.NoError(t, err)
		})
	}
}

func TestConsumerGroupHandler_ConsumeClaim(t *testing.T) {
	value, err := proto.Marshal(&pb.Album{
		Id:     1,
		Title:  "Test Album",
		Artist: "Test Artist",
		Price:  9.99,
	})
	assert.NoError(t, err)

	tests := []struct {
		name           string
		setupMock      func(*MockConsumerGroupSession, *MockConsumerGroupClaim, context.Context, context.CancelFunc)
		expectedErr    bool
		waitForProcess bool
	}{
		{
			name: "consume single message successfully",
			setupMock: func(session *MockConsumerGroupSession, claim *MockConsumerGroupClaim, ctx context.Context, cancel context.CancelFunc) {
				claim.messages = make(chan *sarama.ConsumerMessage, 1)
				claim.messages <- &sarama.ConsumerMessage{
					Value:     []byte(value),
					Topic:     "test-topic",
					Partition: 0,
					Offset:    0,
					Timestamp: time.Now(),
				}

				session.On("Context").Return(ctx)
				session.On("MarkMessage", mock.Anything, "").Run(func(args mock.Arguments) {
					cancel()
				}).Return()
			},
			expectedErr:    false,
			waitForProcess: true,
		},
		{
			name: "handle closed message channel",
			setupMock: func(session *MockConsumerGroupSession, claim *MockConsumerGroupClaim, ctx context.Context, cancel context.CancelFunc) {
				claim.messages = make(chan *sarama.ConsumerMessage)
				close(claim.messages)

				session.On("Context").Return(ctx)
			},
			expectedErr:    false,
			waitForProcess: false,
		},
		{
			name: "handle context cancellation",
			setupMock: func(session *MockConsumerGroupSession, claim *MockConsumerGroupClaim, ctx context.Context, cancel context.CancelFunc) {
				claim.messages = make(chan *sarama.ConsumerMessage)
				cancel()

				session.On("Context").Return(ctx)
			},
			expectedErr:    false,
			waitForProcess: false,
		},
		{
			name: "consume multiple messages",
			setupMock: func(session *MockConsumerGroupSession, claim *MockConsumerGroupClaim, ctx context.Context, cancel context.CancelFunc) {
				claim.messages = make(chan *sarama.ConsumerMessage, 3)
				claim.messages <- &sarama.ConsumerMessage{
					Value:     []byte(value),
					Topic:     "test-topic",
					Partition: 0,
					Offset:    0,
					Timestamp: time.Now(),
				}
				claim.messages <- &sarama.ConsumerMessage{
					Value:     []byte(value),
					Topic:     "test-topic",
					Partition: 0,
					Offset:    1,
					Timestamp: time.Now(),
				}
				claim.messages <- &sarama.ConsumerMessage{
					Value:     []byte(value),
					Topic:     "test-topic",
					Partition: 0,
					Offset:    2,
					Timestamp: time.Now(),
				}

				callCount := 0
				session.On("Context").Return(ctx)
				session.On("MarkMessage", mock.Anything, "").Run(func(args mock.Arguments) {
					callCount++
					if callCount == 3 {
						cancel()
					}
				}).Return()
			},
			expectedErr:    false,
			waitForProcess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cgh := &consumerGroupHandler{
				Ready: make(chan bool),
			}

			mockSession := new(MockConsumerGroupSession)
			mockClaim := new(MockConsumerGroupClaim)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setupMock(mockSession, mockClaim, ctx, cancel)

			err := cgh.ConsumeClaim(mockSession, mockClaim)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockSession.AssertExpectations(t)
		})
	}
}

func TestConsumerGroupHandler_NewInstance(t *testing.T) {
	cgh := &consumerGroupHandler{
		Ready: make(chan bool),
	}

	assert.NotNil(t, cgh)
	assert.NotNil(t, cgh.Ready)
}
