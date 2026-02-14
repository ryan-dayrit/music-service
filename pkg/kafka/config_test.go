package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Struct(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   Config
	}{
		{
			name: "default config",
			config: Config{
				Brokers:       "localhost:9092",
				Topics:        "test-topic",
				ConsumerGroup: "test-group",
				Assignor:      "range",
				Oldest:        false,
			},
			want: Config{
				Brokers:       "localhost:9092",
				Topics:        "test-topic",
				ConsumerGroup: "test-group",
				Assignor:      "range",
				Oldest:        false,
			},
		},
		{
			name: "config with multiple brokers",
			config: Config{
				Brokers:       "localhost:9092,localhost:9093,localhost:9094",
				Topics:        "topic1,topic2",
				ConsumerGroup: "my-group",
				Assignor:      "sticky",
				Oldest:        true,
			},
			want: Config{
				Brokers:       "localhost:9092,localhost:9093,localhost:9094",
				Topics:        "topic1,topic2",
				ConsumerGroup: "my-group",
				Assignor:      "sticky",
				Oldest:        true,
			},
		},
		{
			name: "empty config",
			config: Config{
				Brokers:       "",
				Topics:        "",
				ConsumerGroup: "",
				Assignor:      "",
				Oldest:        false,
			},
			want: Config{
				Brokers:       "",
				Topics:        "",
				ConsumerGroup: "",
				Assignor:      "",
				Oldest:        false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want.Brokers, tt.config.Brokers)
			assert.Equal(t, tt.want.Topics, tt.config.Topics)
			assert.Equal(t, tt.want.ConsumerGroup, tt.config.ConsumerGroup)
			assert.Equal(t, tt.want.Assignor, tt.config.Assignor)
			assert.Equal(t, tt.want.Oldest, tt.config.Oldest)
		})
	}
}

func TestConfig_FieldTags(t *testing.T) {
	config := Config{
		Brokers:       "broker1:9092",
		Topics:        "topic1",
		ConsumerGroup: "group1",
		Assignor:      "roundrobin",
		Oldest:        true,
	}

	assert.Equal(t, "broker1:9092", config.Brokers)
	assert.Equal(t, "topic1", config.Topics)
	assert.Equal(t, "group1", config.ConsumerGroup)
	assert.Equal(t, "roundrobin", config.Assignor)
	assert.True(t, config.Oldest)
}
