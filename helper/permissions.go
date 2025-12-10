package helper

import "github.com/gofiber/fiber/v2"

func HasPermission(c *fiber.Ctx, requiredPerm string) bool {
    perms, ok := c.Locals("permissions").([]string)
    if !ok {
        return false
    }

    for _, p := range perms {
        if p == requiredPerm {
            return true
        }
    }
    
    return false
}