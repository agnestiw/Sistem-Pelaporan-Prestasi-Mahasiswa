package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// NewFiberApp menginisialisasi Fiber dengan konfigurasi custom
func NewFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
		// Nama aplikasi di header response
		AppName: "Sistem Pelaporan Prestasi Mahasiswa v1.0",
		
		// Batas ukuran body request (penting untuk upload file bukti prestasi)
		BodyLimit: 10 * 1024 * 1024, // 10 MB
		
		// Custom Error Handler agar response selalu JSON (tidak HTML default)
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		},
	})

	// Tambahkan middleware Logger bawaan Fiber
	// Ini akan mencatat setiap request HTTP (GET, POST, dll) ke console/log
	app.Use(logger.New())

	return app
}