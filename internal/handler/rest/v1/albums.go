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

// ErrorResponse describes an API error payload.
type ErrorResponse struct {
	Error string `json:"error" example:"cannot parse JSON"`
}

func NewAlbumsHandler(producerHandler kafka.ProducerHandler, repository orm.Repository) *albumsHandler {
	return &albumsHandler{
		producerHandler: producerHandler,
		repository:      repository,
	}
}

// @Summary Create albums
// @Description Queues multiple albums for asynchronous persistence.
// @ID create-albums
// @Tags albums
// @Accept json
// @Produce json
// @Param albums body []pb.Album true "Album payload"
// @Success 201 {array} pb.Album
// @Failure 400 {object} ErrorResponse
// @Router /albums [post]
func (h *albumsHandler) CreateAlbums(ctx *fiber.Ctx) error {
	return h.handleAlbumsWrite(ctx)
}

// @Summary Upsert albums
// @Description Queues multiple albums for asynchronous persistence.
// @ID upsert-albums
// @Tags albums
// @Accept json
// @Produce json
// @Param albums body []pb.Album true "Album payload"
// @Success 201 {array} pb.Album
// @Failure 400 {object} ErrorResponse
// @Router /albums [put]
func (h *albumsHandler) UpsertAlbums(ctx *fiber.Ctx) error {
	return h.handleAlbumsWrite(ctx)
}

func (h *albumsHandler) handleAlbumsWrite(ctx *fiber.Ctx) error {
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

// @Summary Get all albums
// @Description Returns all persisted albums.
// @ID get-albums
// @Tags albums
// @Produce json
// @Success 200 {array} pb.Album
// @Failure 500 {object} ErrorResponse
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
