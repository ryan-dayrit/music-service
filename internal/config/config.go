package config

import (
	"music-service/pkg/grpc"
	"music-service/pkg/kafka"
	"music-service/pkg/postgres/db"
)

type Config struct {
	Grpc     grpc.Config  `yaml:"grpc"`
	Database db.Config    `yaml:"database"`
	Consumer kafka.Config `yaml:"consumer"`
}
