package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"postgres-crud/config"
	"postgres-crud/dal/album"
	"postgres-crud/db"
)

var directCmd = &cobra.Command{
	Use:   "direct",
	Short: "shows the albums",
	Long:  `scans the postgres db directory and shows the albums`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("Error loading config %v", err)
		}

		db, err := db.GetDB(*cfg)
		if err != nil {
			log.Fatalf("Error getting db: %v", err)
		}
		defer db.Close()

		repository := album.NewRepository(db)

		albums, err := repository.Read()
		if err != nil {
			log.Fatalf("Error reading albums: %v", err)
		}

		for _, v := range albums {
			log.Println(v)
		}
	},
}
