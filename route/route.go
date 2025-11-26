package route

import (
	"sistem-prestasi/app/service" 
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authService *service.AuthService) {
	
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Server berjalan dengan koneksi Hybrid (Postgres & Mongo)",
		})
	})

	api := app.Group("/api/v1")

	api.Post("/auth/login", authService.Login)
}