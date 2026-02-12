package grpc

type Config struct {
	Network string `yaml:"network"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}
