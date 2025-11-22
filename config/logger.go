package config

import (
	"log"
	"os"
)

// SetupLogger mengatur output log ke file dan console
func SetupLogger() {
	// 1. Pastikan folder logs ada
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		_ = os.Mkdir("logs", 0755)
	}

	// 2. Buka atau Buat file app.log
	// os.O_APPEND: Menambahkan teks di akhir file (tidak menimpa)
	// os.O_CREATE: Membuat file jika belum ada
	// os.O_WRONLY: Mode write only
	file, err := os.OpenFile("logs/app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	// 3. Set output log default Go ke file tersebut
	// Jadi setiap kali Anda panggil log.Println(), akan masuk ke app.log
	log.SetOutput(file)
	
	// Opsional: Jika ingin log tampil di console JUGA, gunakan MultiWriter (perlu import "io")
	// tapi untuk sekarang kita simpan ke file saja agar rapi sesuai screenshot.
}