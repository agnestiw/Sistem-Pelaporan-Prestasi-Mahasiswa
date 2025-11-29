package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"sistem-prestasi/config"
	"sistem-prestasi/database"
)

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