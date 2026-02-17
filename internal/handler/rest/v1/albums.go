package v1

import (
	"github.com/gofiber/fiber/v2"

	"music-service/gen/pb"

	"music-service/internal/repository/postgres/orm"
	"music-service/pkg/kafka"
)

type albumsHandler struct {
	producerHandler kafka.ProducerHandler
	repository      orm.Repository
}

func NewAlbumsHandler(producerHandler kafka.ProducerHandler, repository orm.Repository) *albumsHandler {
	return &albumsHandler{
		producerHandler: producerHandler,
		repository:      repository,
	}
}

// @Summary Creates albums
// @ID create-albums
// @Produce json
// @Success 201 {array} pb.Album
// @Router /albums [post] [put]
func (h *albumsHandler) CreateAlbums(ctx *fiber.Ctx) error {
	newAlbums := []*pb.Album{}
	if err := ctx.BodyParser(&newAlbums); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	for _, newAlbum := range newAlbums {
		h.producerHandler.Produce(ctx.Context(), newAlbum)
	}
	return ctx.Status(fiber.StatusCreated).JSON(newAlbums)
}

// @Summary Gets all albums
// @ID get-albums
// @Produce json
// @Success 200 {array} pb.Album
// @Router /albums [get]
func (h *albumsHandler) GetAlbums(ctx *fiber.Ctx) error {
	albums, err := h.repository.Get()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get albums",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(albums)
}
