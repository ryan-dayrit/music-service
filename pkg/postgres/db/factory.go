package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	dataSourceNameTemplate := "user=%s dbname=%s sslmode=%s password=%s host=%s"
	dataSourceName :=
		fmt.Sprintf(
			dataSourceNameTemplate,
			cfg.User,
			cfg.DBName,
			cfg.SSLMode,
			cfg.Password,
			cfg.Host,
		)
	db, err := sqlx.Connect(cfg.DriverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
