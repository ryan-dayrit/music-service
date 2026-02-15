package postgres

import (
	"log"

	"github.com/spf13/cobra"

	"music-service/internal/config"
	"music-service/internal/repository/postgres/sqlx"
	"music-service/pkg/postgres/sqlx/db"
)

func NewPostgresGetAllCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "postgres-get-all",
		Short: "shows the albums",
		Long:  `queries the postgres db directly and shows the albums`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				log.Fatalf("failed to load config %v", err)
			}

			db, err := db.NewDB(cfg.Postgres)
			if err != nil {
				log.Fatalf("failed to get db: %v", err)
			}
			defer db.Close()

			repository := sqlx.NewRepository(db)

			albums, err := repository.Read()
			if err != nil {
				log.Fatalf("failed to read albums: %v", err)
			}

			for _, v := range albums {
				log.Println(v)
			}
		},
	}
}
