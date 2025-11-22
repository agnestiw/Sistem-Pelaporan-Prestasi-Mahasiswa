package main

import (
	"log"

	"sistem-prestasi/config"   // Import folder config
	"sistem-prestasi/database" // Import folder database

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Setup Logger (Membuat file app.log otomatis)
	config.SetupLogger()
	log.Println("Starting Application...") // Ini akan tertulis di logs/app.log

	// 2. Load Config Environment
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	// 3. Connect Database
	database.Connect()

	// 4. Init Fiber App (dari config/app.go)
	app := config.NewFiberApp()

	// Route Test
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Sistem Pelaporan Prestasi API is running",
			"docs":    "See SRS for details",
		})
	})

	// 5. Start Server
	log.Printf("Server running on port %s", cfg.AppPort)
	if err := app.Listen(cfg.AppPort); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}