package v1

import (
	"github.com/gofiber/fiber/v2"

	"music-service/gen/pb"
	"music-service/internal/handler/kafka"
	"music-service/internal/repository/postgres/orm"
)

type albumsHandler struct {
	producer   kafka.ProducerHandler
	repository orm.Repository
}

func NewAlbumsHandler(producer kafka.ProducerHandler, repository orm.Repository) *albumsHandler {
	return &albumsHandler{
		producer:   producer,
		repository: repository,
	}
}

func (h *albumsHandler) CreateAlbums(ctx *fiber.Ctx) error {
	newAlbums := []*pb.Album{}
	if err := ctx.BodyParser(&newAlbums); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	for _, newAlbum := range newAlbums {
		h.producer.Produce(ctx.Context(), newAlbum)
	}
	return ctx.Status(fiber.StatusCreated).JSON(newAlbums)
}

func (h *albumsHandler) GetAlbums(ctx *fiber.Ctx) error {
	albums, err := h.repository.Get()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get albums",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(albums)
}
