package postgres

import (
	"log"
	"math/rand/v2"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"

	"music-service/internal/config"
	"music-service/internal/models"
	"music-service/internal/repository/postgres/orm"
	"music-service/pkg/postgres/orm/db"
)

func NewPostgresInsertCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "postgres-insert",
		Short: "inserts a new album",
		Long:  `inserts a new album into the postgres db directly`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				log.Fatalf("failed to load config %v", err)
			}

			db := db.NewDB(cfg.Postgres)
			defer db.Close()

			repository := orm.NewRepository(db)

			album := models.Album{
				Id:     int(rand.Int32()),
				Title:  uuid.NewString(),
				Artist: uuid.NewString(),
				Price:  decimal.NewFromFloat(rand.Float64()),
			}

			repository.Create(album)
			log.Println(album)
		},
	}
}
