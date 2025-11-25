package config

import (
	"log"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Warning: File .env tidak ditemukan, menggunakan environment sistem.")
	}
}