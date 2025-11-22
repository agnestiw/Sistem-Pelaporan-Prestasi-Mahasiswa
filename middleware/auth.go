package middleware

import (
	"strings"

	"sistem-prestasi/config"
	"sistem-prestasi/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protect adalah middleware untuk memvalidasi JWT Token
func Protect(c *fiber.Ctx) error {
	// 1. Ambil Header Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Missing Authorization Header")
	}

	// 2. Cek format "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid Token Format")
	}

	tokenString := parts[1]
	cfg, _ := config.LoadConfig()

	// 3. Parse Token
	token, err := jwt.ParseWithClaims(tokenString, &utils.JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or Expired Token")
	}

	// 4. Simpan Claims ke Context (agar bisa dibaca di handler selanjutnya)
	claims, ok := token.Claims.(*utils.JwtClaims)
	if ok && token.Valid {
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		c.Locals("permissions", claims.Permissions)
	}

	return c.Next()
}