package route

import (
	"sistem-prestasi/app/service"
	"sistem-prestasi/middleware" 
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authService *service.AuthService) {
	
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Server berjalan dengan koneksi Hybrid (Postgres & Mongo)",
		})
	})

	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/login", authService.Login)
	auth.Post("/refresh", authService.Refresh)
	auth.Post("/logout", middleware.Protect(), authService.Logout)
    auth.Get("/profile", middleware.Protect(), authService.Profile)
}