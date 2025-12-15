package route

import (
	"sistem-prestasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	analytics := api.Group("/analytics")
	analytics.Use(middleware.Protect())
}
