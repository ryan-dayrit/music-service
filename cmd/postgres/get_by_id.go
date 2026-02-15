package postgres

import (
	"log"

	"github.com/spf13/cobra"

	"music-service/internal/config"
	"music-service/internal/repository/postgres/orm"
	"music-service/pkg/postgres/orm/db"
)

func NewPostgresGetByIdCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "postgres-get-by-id",
		Short: "shows the album of the given id",
		Long:  `queries the postgres db directly and shows the album of the given id`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				log.Fatalf("failed to load config %v", err)
			}

			db := db.NewDB(cfg.Postgres)
			defer db.Close()

			repository := orm.NewRepository(db)

			album, err := repository.GetById(1)
			if err != nil {
				log.Fatalf("failed to read albums: %v", err)
			}

			log.Println(album)
		},
	}
}
