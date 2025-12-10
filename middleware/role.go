package middleware

import (
    "github.com/gofiber/fiber/v2"
)

func Role(allowed ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {

        userRole := c.Locals("role")
        if userRole == nil {
            return c.Status(401).JSON(fiber.Map{"error": "No role assigned"})
        }

        roleName := userRole.(string)

        for _, v := range allowed {
            if v == roleName {
                return c.Next()
            }
        }

        return c.Status(403).JSON(fiber.Map{"error": "Forbidden: insufficient role"})
    }
}
