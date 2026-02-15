package v1

import (
	"music-service/internal/handler/kafka"
	v1 "music-service/internal/handler/rest/v1"
	"music-service/internal/repository/postgres/orm"

	"github.com/gofiber/fiber/v2"
)

func RegisterPublicRoutes(router fiber.Router, producerHandler kafka.ProducerHandler, repository orm.Repository) {
	albumHandler := v1.NewAlbumHandler(producerHandler)
	router.Post("/album", albumHandler.CreateAlbum)
	router.Put("/album", albumHandler.CreateAlbum)

	albumsHandler := v1.NewAlbumsHandler(producerHandler, repository)
	router.Post("/albums", albumsHandler.CreateAlbums)
	router.Put("/albums", albumsHandler.CreateAlbums)
	router.Get("/albums", albumsHandler.GetAlbums)
}
