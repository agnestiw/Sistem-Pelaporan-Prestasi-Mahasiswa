package config

import (
	"sistem-prestasi/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New(LoggerConfig())) // LoggerConfig dari file logger.go sebelumnya

	// Panggil Setup Routes di sini (Sesuai gaya screenshot Anda)
	route.SetupRoutes(app)

	return app
}