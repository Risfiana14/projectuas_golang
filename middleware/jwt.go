package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"projectuas/app/model"
)

func JWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return fiber.NewError(401, "Missing Authorization header")
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return fiber.NewError(401, "Invalid token")
		}

		claims := token.Claims.(*model.Claims)

		// 🔥 INI YANG KEMARIN HILANG
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

