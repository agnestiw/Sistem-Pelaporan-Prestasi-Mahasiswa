package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"sistem-prestasi/config"
	"sistem-prestasi/database"

	_ "sistem-prestasi/docs"
)

// @title Sistem Prestasi Mahasiswa
// @version 1.0
// @description Ini adalah sistem prestasi mahasiswa
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.LoadEnv()

	database.InitPostgres()
	database.InitMongo()

	defer database.DB.Close()
	defer func() {
		if err := database.MongoClient.Disconnect(context.Background()); err != nil {
			log.Println(err)
		}
	}()

	app := config.NewApp()



	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println("🚀 Server is running on port " + port)
	log.Fatal(app.Listen(":" + port))
}
