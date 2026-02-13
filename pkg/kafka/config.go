package kafka

type Config struct {
	Brokers       string `yaml:"brokers"`
	Topics        string `yaml:"topics"`
	ConsumerGroup string `yaml:"consumer_group"`
}
