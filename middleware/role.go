package middleware

import "github.com/gofiber/fiber/v2"

func Role(allowedRoles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        role := c.Locals("role").(string)

        for _, r := range allowedRoles {
            if role == r {
                return c.Next()
            }
        }

        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "status":  "error",
            "message": "Forbidden: insufficient permissions",
        })
    }
}

