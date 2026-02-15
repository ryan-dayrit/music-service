package db

import (
	"github.com/go-pg/pg/v10"

	"music-service/pkg/postgres"
)

func NewDB(cfg postgres.Config) *pg.DB {
	return pg.Connect(&pg.Options{
		User:     cfg.User,
		Database: cfg.DBName,
		Password: cfg.Password,
		Addr:     cfg.Host + ":5432",
	})
}
