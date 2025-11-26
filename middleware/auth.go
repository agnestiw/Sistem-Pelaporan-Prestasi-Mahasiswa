package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized: Token wajib ada"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized: Format token salah"})
		}
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing tidak valid")
			}
			return []byte(os.Getenv("API_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized: Token tidak valid"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("user_id", claims["user_id"])
			c.Locals("role_id", claims["role_id"])
			
			permInterface := claims["permissions"].([]interface{})
			var permissions []string
			for _, v := range permInterface {
				permissions = append(permissions, v.(string))
			}
			c.Locals("permissions", permissions)
		}

		return c.Next()
	}
}

func HasPermission(requiredPerm string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userPerms, ok := c.Locals("permissions").([]string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Forbidden: Permission tidak ditemukan"})
		}

		for _, p := range userPerms {
			if p == requiredPerm {
				return c.Next() 
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": fmt.Sprintf("Forbidden: Anda tidak memiliki akses '%s'", requiredPerm),
		})
	}
}