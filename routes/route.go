package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func InitRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	InitAuthRoute(api)
	InitUserRoute(api)
	InitPostRoute(api)
}
