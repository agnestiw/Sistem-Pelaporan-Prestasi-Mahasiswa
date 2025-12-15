package route

import (
	"sistem-prestasi/app/service"
	"sistem-prestasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(app *fiber.App) {

	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/login", service.Login)
	auth.Post("/refresh", service.Refresh)
	auth.Post("/logout", middleware.Protect(), service.Logout)
	auth.Get("/profile", middleware.Protect(), service.Profile)

}
