package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/ryan-dayrit/music-service/config"
	"github.com/ryan-dayrit/music-service/dal/album"
	"github.com/ryan-dayrit/music-service/db"
)

var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "shows the albums",
	Long:  `queries the postgres db directly and shows the albums`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("failed to load config %v", err)
		}

		db, err := db.GetDB(*cfg)
		if err != nil {
			log.Fatalf("failed to get db: %v", err)
		}
		defer db.Close()

		repository := album.NewRepository(db)

		albums, err := repository.Read()
		if err != nil {
			log.Fatalf("failed to read albums: %v", err)
		}

		for _, v := range albums {
			log.Println(v)
		}
	},
}
