package helper

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateUniqueFilename membuat nama file unik agar tidak bentrok
// Contoh output: 20251122-uuid-sertifikat.pdf
func GenerateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	filename := strings.TrimSuffix(originalName, ext)
	
	// Bersihkan nama file dari karakter aneh
	filename = strings.ReplaceAll(filename, " ", "_")
	
	uniqueID := uuid.New().String()
	timestamp := time.Now().Format("20060102")

	return fmt.Sprintf("%s-%s-%s%s", timestamp, uniqueID[:8], filename, ext)
}

// IsAllowedFileType mengecek apakah file yang diupload valid (PDF/Gambar)
func IsAllowedFileType(mimeType string) bool {
	allowedTypes := map[string]bool{
		"application/pdf": true,
		"image/jpeg":      true,
		"image/png":       true,
		"image/jpg":       true,
	}
	return allowedTypes[mimeType]
}