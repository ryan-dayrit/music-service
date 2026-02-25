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

// @Summary Create an album
// @Description Queues a single album for asynchronous persistence.
// @ID create-album
// @Tags albums
// @Accept json
// @Produce json
// @Param album body pb.Album true "Album payload"
// @Success 201 {object} pb.Album
// @Failure 400 {object} ErrorResponse
// @Router /album [post]
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
