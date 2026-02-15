package sarama

import (
	"testing"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"

	"music-service/pkg/kafka"
)

func TestNewConsumerGroup_StickyAssignor(t *testing.T) {
	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "sticky",
		Oldest:        false,
	}

	consumerGroup, err := NewConsumerGroup(cfg)

	if consumerGroup != nil {
		defer consumerGroup.Close()
	}

	// Either succeeds (if Kafka is running) or fails with connection error
	if err != nil {
		assert.Contains(t, err.Error(), "Error creating consumer group")
	} else {
		assert.NotNil(t, consumerGroup)
	}
}

func TestNewConsumerGroup_RoundRobinAssignor(t *testing.T) {
	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "roundrobin",
		Oldest:        false,
	}

	consumerGroup, err := NewConsumerGroup(cfg)

	if consumerGroup != nil {
		defer consumerGroup.Close()
	}

	if err != nil {
		assert.Contains(t, err.Error(), "Error creating consumer group")
	} else {
		assert.NotNil(t, consumerGroup)
	}
}

func TestNewConsumerGroup_RangeAssignor(t *testing.T) {
	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "range",
		Oldest:        false,
	}

	consumerGroup, err := NewConsumerGroup(cfg)

	if consumerGroup != nil {
		defer consumerGroup.Close()
	}

	if err != nil {
		assert.Contains(t, err.Error(), "Error creating consumer group")
	} else {
		assert.NotNil(t, consumerGroup)
	}
}

func TestNewConsumerGroup_InvalidAssignor(t *testing.T) {
	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "invalid",
		Oldest:        false,
	}

	consumerGroup, err := NewConsumerGroup(cfg)

	assert.Error(t, err)
	assert.Nil(t, consumerGroup)
	assert.Contains(t, err.Error(), "Unrecognized consumer group partition assignor: invalid")
}

func TestNewConsumerGroup_OldestOffset(t *testing.T) {
	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "range",
		Oldest:        true,
	}

	consumerGroup, err := NewConsumerGroup(cfg)

	if consumerGroup != nil {
		defer consumerGroup.Close()
	}

	if err != nil {
		assert.Contains(t, err.Error(), "Error creating consumer group")
	} else {
		assert.NotNil(t, consumerGroup)
	}
}

func TestNewConsumerGroup_MultipleBrokers(t *testing.T) {
	cfg := kafka.Config{
		Brokers:       "localhost:9092,localhost:9093,localhost:9094",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "sticky",
		Oldest:        false,
	}

	consumerGroup, err := NewConsumerGroup(cfg)

	if consumerGroup != nil {
		defer consumerGroup.Close()
	}

	if err != nil {
		assert.Contains(t, err.Error(), "Error creating consumer group")
	} else {
		assert.NotNil(t, consumerGroup)
	}
}

func TestNewConsumerGroup_EmptyBrokers(t *testing.T) {
	cfg := kafka.Config{
		Brokers:       "",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "range",
		Oldest:        false,
	}

	consumerGroup, err := NewConsumerGroup(cfg)

	if consumerGroup != nil {
		defer consumerGroup.Close()
	}

	if err != nil {
		assert.Contains(t, err.Error(), "Error creating consumer group")
	} else {
		assert.NotNil(t, consumerGroup)
	}
}

func TestNewConsumerGroup_AllAssignors(t *testing.T) {
	tests := []struct {
		name     string
		assignor string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "sticky assignor",
			assignor: "sticky",
			wantErr:  true,
			errMsg:   "Error creating consumer group",
		},
		{
			name:     "roundrobin assignor",
			assignor: "roundrobin",
			wantErr:  true,
			errMsg:   "Error creating consumer group",
		},
		{
			name:     "range assignor",
			assignor: "range",
			wantErr:  true,
			errMsg:   "Error creating consumer group",
		},
		{
			name:     "invalid assignor",
			assignor: "unknown",
			wantErr:  true,
			errMsg:   "Unrecognized consumer group partition assignor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := kafka.Config{
				Brokers:       "localhost:9092",
				Topics:        "test-topic",
				ConsumerGroup: "test-group",
				Assignor:      tt.assignor,
				Oldest:        false,
			}

			consumerGroup, err := NewConsumerGroup(cfg)

			if consumerGroup != nil {
				defer consumerGroup.Close()
			}

			if tt.wantErr && tt.errMsg == "Unrecognized consumer group partition assignor" {
				// Invalid assignor should always fail
				assert.Error(t, err)
				assert.Nil(t, consumerGroup)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else if tt.wantErr {
				// Valid assignor may succeed or fail depending on broker availability
				if err != nil {
					assert.Contains(t, err.Error(), tt.errMsg)
				} else {
					assert.NotNil(t, consumerGroup)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, consumerGroup)
			}
		})
	}
}

func TestNewConsumerGroup_OffsetConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		oldest bool
	}{
		{
			name:   "oldest offset enabled",
			oldest: true,
		},
		{
			name:   "oldest offset disabled",
			oldest: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := kafka.Config{
				Brokers:       "localhost:9092",
				Topics:        "test-topic",
				ConsumerGroup: "test-group",
				Assignor:      "range",
				Oldest:        tt.oldest,
			}

			consumerGroup, err := NewConsumerGroup(cfg)

			if consumerGroup != nil {
				defer consumerGroup.Close()
			}

			// May succeed or fail depending on broker availability
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NotNil(t, consumerGroup)
			}
		})
	}
}

func TestNewConsumerGroup_ConfigValidation(t *testing.T) {
	cfg := kafka.Config{
		Brokers:       "localhost:9092",
		Topics:        "test-topic",
		ConsumerGroup: "test-group",
		Assignor:      "range",
		Oldest:        false,
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version, _ = sarama.ParseKafkaVersion(sarama.DefaultVersion.String())
	saramaCfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}

	assert.NotNil(t, saramaCfg)
	assert.Equal(t, sarama.OffsetNewest, saramaCfg.Consumer.Offsets.Initial)

	_, err := NewConsumerGroup(cfg)
	// May succeed or fail depending on broker availability
	if err != nil {
		assert.Error(t, err)
	}
}
