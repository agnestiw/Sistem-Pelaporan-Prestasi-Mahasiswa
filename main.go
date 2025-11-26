package main

import (
	"log"
	"os"

	"sistem-prestasi/config"
	"sistem-prestasi/database"
	"sistem-prestasi/route"
	
	"sistem-prestasi/app/repository/postgre"
	"sistem-prestasi/app/service"
)

func main() {
	config.LoadEnv()
	dbPostgres := database.InitPostgres()
	// dbMongo := database.InitMongo() 
	defer dbPostgres.Close()

	userRepo := postgre.NewUserRepository(dbPostgres)
	
	authService := service.NewAuthService(userRepo)

	app := config.NewApp()
	
	route.SetupRoutes(app, authService)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}