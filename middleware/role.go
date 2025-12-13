package middleware

import "github.com/gofiber/fiber/v2"

func Role(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		roleI := c.Locals("role")
		if roleI == nil {
			return fiber.NewError(401, "Role not found")
		}

		role, ok := roleI.(string)
		if !ok {
			return fiber.NewError(401, "Invalid role")
		}

		for _, allowed := range allowedRoles {
			if role == allowed {
				return c.Next()
			}
		}

		return fiber.NewError(403, "Forbidden: insufficient role")
	}
}
