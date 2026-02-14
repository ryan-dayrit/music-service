package config

import (
	"music-service/pkg/grpc"
	"music-service/pkg/kafka"
	"music-service/pkg/postgres/db"
	"music-service/pkg/rest"
)

type Config struct {
	Grpc     grpc.Config  `yaml:"grpc"`
	Database db.Config    `yaml:"database"`
	Kafka    kafka.Config `yaml:"kafka"`
	Rest     rest.Config  `yaml:"rest"`
}
