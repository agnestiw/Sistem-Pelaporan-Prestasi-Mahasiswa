package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewApp() *fiber.App {
	app := fiber.New()

	// Middleware
	app.Use(cors.New())
	
	// Panggil konfigurasi logger dari file logger.go
	app.Use(logger.New(LoggerConfig()))

	// Route tidak didaftarkan di sini. App hanya menyiapkan instance.
	return app
}