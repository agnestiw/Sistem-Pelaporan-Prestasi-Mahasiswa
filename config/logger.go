package config

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

func LoggerConfig() logger.Config {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		_ = os.Mkdir("logs", 0755)
	}

	fileName := "logs/app.log"

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