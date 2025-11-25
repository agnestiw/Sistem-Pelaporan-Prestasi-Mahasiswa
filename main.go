package main

import (
	"log"
	"os"

	"sistem-prestasi/config"
	"sistem-prestasi/database"
	"sistem-prestasi/route" // Import paket route

	// Nanti kita akan import repository dan service di sini
)

func main() {
	// 1. Load Environment Variables
	config.LoadEnv()

	// 2. Inisialisasi Database
	dbPostgres := database.InitPostgres()
	database.InitMongo()

	defer dbPostgres.Close()

	// 3. Setup Fiber App
	app := config.NewApp()

	// ---------------------------------------------------------
	// AREA DEPENDENCY INJECTION 
	// ---------------------------------------------------------

	// 4. Wiring Routes
	route.SetupRoutes(app)

	// 5. Jalankan Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	
	log.Fatal(app.Listen(":" + port))
}