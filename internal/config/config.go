package config

import (
	"music-service/pkg/postgres/db"
)

type Config struct {
	Service struct {
		Network string `yaml:"network"`
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
	} `yaml:"service"`

	Database db.Config `yaml:"database"`
}
