package config

import (
	"music-service/pkg/grpc"
	"music-service/pkg/kafka"
	"music-service/pkg/postgres"
	"music-service/pkg/rest"
)

type Config struct {
	Grpc     grpc.Config     `yaml:"grpc"`
	Postgres postgres.Config `yaml:"postgres"`
	Kafka    kafka.Config    `yaml:"kafka"`
	Rest     rest.Config     `yaml:"rest"`
}
