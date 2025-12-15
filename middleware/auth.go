package middleware

import (
	"strings"

	"sistem-prestasi/helper"
	memory "sistem-prestasi/memory"

	"github.com/gofiber/fiber/v2"
)

func Protect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized: Token wajib ada",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized: Format token salah",
			})
		}

		tokenString := parts[1]

		if memory.IsBlacklisted(tokenString) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized: Anda sudah logout, silakan login kembali",
			})
		}

		claims, err := helper.ValidateJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized: Token tidak valid",
			})
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("role_id", claims["role_id"])
		c.Locals("role_name", claims["role_name"])
		c.Locals("student_id", claims["student_id"])
		c.Locals("advisor_id", claims["advisor_id"])

		var permissions []string
		if permClaim, ok := claims["permissions"]; ok && permClaim != nil {
			if permInterface, ok := permClaim.([]interface{}); ok {
				for _, v := range permInterface {
					if s, ok := v.(string); ok {
						permissions = append(permissions, s)
					}
				}
			}
		}
		c.Locals("permissions", permissions)

		return c.Next()
	}
}


