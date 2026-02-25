package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// HealthResponse describes the health check response body.
type HealthResponse struct {
	Status string `json:"status" example:"OK"`
}

// @Summary Health check
// @Description Returns service health status.
// @ID health-check
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func RegisterHealthRoute(router fiber.Router) {
	router.Get(
		"/health",
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"status": http.StatusText(http.StatusOK),
			})
		},
	)
}
