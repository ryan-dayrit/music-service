package v1

import (
	"music-service/internal/handler/kafka"
	v1 "music-service/internal/handler/rest/v1"

	"github.com/gofiber/fiber/v2"
)

func RegisterPublicRoutes(r fiber.Router, producer kafka.Producer) {
	albumHandler := v1.NewAlbumHandler(producer)
	r.Post("/albums", albumHandler.CreateAlbum)
}
