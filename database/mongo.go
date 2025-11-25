package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo() *mongo.Database {
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Gagal membuat client Mongo: %v", err)
	}

	// Cek koneksi (Ping)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Gagal terhubung ke MongoDB: %v", err)
	}

	log.Println("✅ Berhasil terhubung ke MongoDB")
	
	// Mengembalikan instance Database agar bisa langsung akses collection nanti
	return client.Database(dbName) 
}