package db

import (
	"fmt"

	"postgres-crud/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetDB(cfg config.Config) (*sqlx.DB, error) {
	dataSourceNameTemplate := "user=%s dbname=%s sslmode=%s password=%s host=%s"
	dataSourceName :=
		fmt.Sprintf(
			dataSourceNameTemplate,
			cfg.Database.User,
			cfg.Database.DBName,
			cfg.Database.SSLMode,
			cfg.Database.Password,
			cfg.Database.Host,
		)
	db, err := sqlx.Connect(cfg.Database.DriverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
