package routes

import "github.com/gofiber/fiber/v2"

func RegisterNotFoundRoute(a *fiber.App) {
	a.Use(
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "endpoint is not found",
			})
		},
	)
}
