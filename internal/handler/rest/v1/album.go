package v1

import (
	"github.com/gofiber/fiber/v2"

	"music-service/gen/pb"
	"music-service/pkg/kafka"
)

type albumHandler struct {
	producerHandler kafka.ProducerHandler
}

func NewAlbumHandler(producerHandler kafka.ProducerHandler) *albumHandler {
	return &albumHandler{
		producerHandler: producerHandler,
	}
}

func (h *albumHandler) CreateAlbum(ctx *fiber.Ctx) error {
	newAlbum := &pb.Album{}
	if err := ctx.BodyParser(newAlbum); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	h.producerHandler.Produce(ctx.Context(), newAlbum)
	return ctx.Status(fiber.StatusCreated).JSON(newAlbum)
}
