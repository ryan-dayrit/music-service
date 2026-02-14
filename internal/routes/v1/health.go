package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

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
