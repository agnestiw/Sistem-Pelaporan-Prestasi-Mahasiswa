package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JwtClaims mendefinisikan payload token sesuai kebutuhan RBAC SRS
type JwtClaims struct {
	UserID      string   `json:"user_id"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"` // Disimpan di token agar tidak query DB terus menerus (FR-002)
	jwt.RegisteredClaims
}

// GenerateToken membuat token baru saat login
func GenerateToken(userID, role string, permissions []string, secret string) (string, error) {
	claims := JwtClaims{
		UserID:      userID,
		Role:        role,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid 24 jam
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}