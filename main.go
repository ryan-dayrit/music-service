package main

import (
	"log"

	"postgres-crud/config"
	"postgres-crud/dal/album"
	"postgres-crud/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config %v", err)
	}

	db, err := db.GetDB(*cfg)
	if err != nil {
		log.Fatalf("Error getting db: %v", err)
	}
	defer db.Close()

	repository := album.NewAlbumRepository(db)

	albums, err := repository.Read()
	if err != nil {
		log.Fatalf("Error reading albums: %v", err)
	}

	for _, v := range albums {
		log.Println(v)
	}

}
