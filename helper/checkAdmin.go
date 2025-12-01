package helper

import "github.com/gofiber/fiber/v2"

func IsAdmin(c *fiber.Ctx) bool {
	perms, ok := c.Locals("permissions").([]string)
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == "user:manage" {
			return true
		}
	}
	return false
}
