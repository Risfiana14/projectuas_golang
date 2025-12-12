package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// Role accepts one or more allowed role names.
func Role(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		r := c.Locals("role")
		role, _ := r.(string) // safe cast: if nil/other -> ""
		// if no role in context → unauthenticated (but JWT middleware normally sets it)
		if role == "" {
			return c.Status(403).JSON(fiber.Map{"error": "Forbidden: missing role"})
		}

		for _, a := range allowedRoles {
			if role == a {
				// permitted
				return c.Next()
			}
		}
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden: insufficient role"})
	}
}
