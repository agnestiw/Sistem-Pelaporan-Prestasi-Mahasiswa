package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global Variable
var MongoClient *mongo.Client
var MongoDb *mongo.Database

func InitMongo() {
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Gagal membuat client Mongo: %v", err)
	}

	if err := MongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("Gagal terhubung ke MongoDB: %v", err)
	}

	MongoDb = MongoClient.Database(dbName)
	log.Println("✅ Berhasil terhubung ke MongoDB")
}