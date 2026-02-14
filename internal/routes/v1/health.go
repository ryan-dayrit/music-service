package v1

import "github.com/gofiber/fiber/v2"

func RegisterHealthRoute(r fiber.Router) {
	r.Get(
		"/health",
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"status": "ok",
			})
		},
	)
}
