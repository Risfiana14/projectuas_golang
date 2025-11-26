package middleware

import (
    "os"
    "strings"
    "projectuas/app/model"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    
)

func JWT() fiber.Handler {
    return func(c *fiber.Ctx) error {
        auth := c.Get("Authorization")
        if !strings.HasPrefix(auth, "Bearer ") {
            return c.Status(401).JSON(fiber.Map{"message": "Token diperlukan"})
        }
        tokenStr := strings.TrimPrefix(auth, "Bearer ")
        claims := &model.Claims{}

        token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        }, jwt.WithValidMethods([]string{"HS256"}))

        if err != nil || !token.Valid {
            return c.Status(401).JSON(fiber.Map{"message": "Token tidak valid"})
        }

        c.Locals("user_id", claims.UserID)
        c.Locals("role", claims.Role)
        return c.Next()
    }
}

func Role(allowed ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        role := c.Locals("role").(string)
        for _, r := range allowed {
            if r == role {
                return c.Next()
            }
        }
        return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak"})
    }
}