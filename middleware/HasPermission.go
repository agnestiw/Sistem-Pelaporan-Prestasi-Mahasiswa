package middleware

import (

	"github.com/gofiber/fiber/v2"
)

func HasPermission(requiredPerm string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		userPerms, ok := c.Locals("permissions").([]string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Forbidden: Permission tidak ditemukan",
			})
		}

		for _, p := range userPerms {
			if p == requiredPerm {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"message": "maaf, anda tidak memiliki akses",
			"permission": userPerms,
		})
	}
}
