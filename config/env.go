package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config struct menyimpan semua konfigurasi aplikasi
type Config struct {
	AppPort      string
	DatabaseURL  string // DSN Postgres
	MongoURI     string
	MongoDBName  string
	JWTSecret    string
}

// LoadConfig memuat .env dan mengembalikan struct Config
func LoadConfig() (*Config, error) {
	// Load file .env jika ada
	// Abaikan error jika file tidak ada (berguna saat deployment production yang menggunakan env sistem)
	_ = godotenv.Load()

	return &Config{
		AppPort:     getEnv("APP_PORT", ":3000"),
		DatabaseURL: os.Getenv("DATABASE_URL"), // Nanti kita rakit di main/db agar dinamis
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName: getEnv("MONGO_DB_NAME", "sistem_prestasi"),
		JWTSecret:   getEnv("JWT_SECRET", "secret_default"),
	}, nil
}

// Helper function untuk membaca env dengan default value (fallback)
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}