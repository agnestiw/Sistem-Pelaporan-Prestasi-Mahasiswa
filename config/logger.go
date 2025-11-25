package config

import (
	"fmt"
	"os"
	// "time" // Kita tidak butuh package time lagi untuk nama file

	"github.com/gofiber/fiber/v2/middleware/logger"
)

func LoggerConfig() logger.Config {
	// 1. Pastikan folder logs ada
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		_ = os.Mkdir("logs", 0755)
	}

	// 2. Tentukan nama file log statis (TIDAK PAKAI TANGGAL)
	fileName := "logs/app.log"

	// 3. Buka file dengan mode APPEND
	// os.O_APPEND sangat penting di sini agar log lama tidak tertimpa saat server restart
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("⚠️ Gagal membuka file log, output akan dialihkan ke terminal:", err)
		return logger.Config{
			Format: "[${time}] ${status} - ${method} ${path}\n",
		}
	}

	return logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} - ${latency}\n",
		TimeFormat: "15:04:05",
		TimeZone:   "Asia/Jakarta",
		Output:     file,
	}
}