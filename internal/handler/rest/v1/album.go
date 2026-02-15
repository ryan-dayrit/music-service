package v1

import (
	"github.com/gofiber/fiber/v2"

	"music-service/gen/pb"
	"music-service/internal/handler/kafka"
)

type albumHandler struct {
	producer kafka.ProducerHandler
}

func NewAlbumHandler(producer kafka.ProducerHandler) *albumHandler {
	return &albumHandler{
		producer: producer,
	}
}

func (h *albumHandler) CreateAlbum(ctx *fiber.Ctx) error {
	newAlbum := &pb.Album{}
	if err := ctx.BodyParser(newAlbum); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	h.producer.Produce(ctx.Context(), newAlbum)
	return ctx.Status(fiber.StatusCreated).JSON(newAlbum)
}
