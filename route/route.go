package route

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes mendaftarkan semua route API
// Kita menerima instance Fiber App sebagai argumen.
func SetupRoutes(app *fiber.App) {
	// Root Endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Server berjalan dengan koneksi Hybrid (Postgres & Mongo)",
		})
	})

	// Grouping untuk versi API (sesuai SRS Bab 5.1: /api/v1)
	api := app.Group("/api/v1")

	// -------------------------------------------------------------------
	// AREA ROUTE MODULES (akan diisi di fase berikutnya, misal Auth/Achievement)
	// Contoh: AuthRoutes(api)
	// -------------------------------------------------------------------

	// Contoh rute Auth (diperlukan untuk testing nanti)
	api.Post("/auth/login", func(c *fiber.Ctx) error {
		// Placeholder untuk Auth Service
		return c.JSON(fiber.Map{"message": "Login Endpoint Placeholder"})
	})
}