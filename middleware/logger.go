package middleware

import "github.com/gofiber/fiber/v2"

func Logger() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Simple logger (bisa diganti fiber.Logger nanti)
        return c.Next()
    }
}