// route/route.go
package route

import (
	"time"
	"projectuas/app/model"
	"projectuas/app/service"
	"projectuas/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(app *fiber.App, pgDB *sql.DB, mongoClient *mongo.Client) {
	// LOGIN (Public
	app.Post("/login", func(c *fiber.Ctx) error {
		type Req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var r Req
		if err := c.BodyParser(&r); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "Format salah"})
		}

		users := map[string]string{
			"admin":     "admin123",
			"mahasiswa": "mhs123",
			"dosen":     "dosen123",
		}

		if pass, ok := users[r.Username]; ok && pass == r.Password {
			role := r.Username
			if role == "dosen" {
				role = "dosen_wali"
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Claims{
				UserID: uuid.New(),
				Role:   role,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				},
			})

			signed, _ := token.SignedString([]byte("superrahasia123456789"))
			return c.JSON(fiber.Map{
				"message":      "Login berhasil!",
				"access_token": signed,
				"role":         role,
			})
		}
		return c.Status(401).JSON(fiber.Map{"message": "Username/password salah"})
	})

	// API dengan JWT
	api := app.Group("/api", middleware.JWT())

	// Hanya mahasiswa & admin yang boleh hapus draft
	api.Delete("/achievement/:id", middleware.Role("mahasiswa", "admin"), func(c *fiber.Ctx) error {
		return service.DeleteAchievement(c, pgDB, mongoClient)
	})

	// Test token
	api.Get("/me", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Token valid",
			"user_id": c.Locals("user_id"),
			"role":    c.Locals("role"),
		})
	})
}