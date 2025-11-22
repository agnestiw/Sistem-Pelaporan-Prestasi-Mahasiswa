package middleware

import (
	"sistem-prestasi/utils"

	"github.com/gofiber/fiber/v2"
)

// RequirePermission memastikan user memiliki permission tertentu
func RequirePermission(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil permissions user dari Locals (diset oleh middleware Protect)
		userPermissions, ok := c.Locals("permissions").([]string)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "No permissions found")
		}

		// Cek apakah permission yang dibutuhkan ada di list userPermissions
		for _, p := range userPermissions {
			if p == requiredPermission {
				return c.Next() // Boleh lanjut
			}
		}

		return utils.ErrorResponse(c, fiber.StatusForbidden, "Insufficient Permissions: "+requiredPermission)
	}
}