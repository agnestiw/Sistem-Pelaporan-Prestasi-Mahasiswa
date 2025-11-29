package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Global Variable agar bisa diakses dari package lain
var DB *sql.DB

func InitPostgres() {
	var err error
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi Postgres: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Gagal terhubung ke Postgres: %v", err)
	}

	log.Println("✅ Berhasil terhubung ke PostgreSQL")
}