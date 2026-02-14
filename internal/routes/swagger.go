package routes

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

func RegisterSwaggerRoute(app *fiber.App) {
	route := app.Group("/swagger")
	route.Get("*", swagger.HandlerDefault)
}
