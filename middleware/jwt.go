package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"projectuas/app/repository"
)

func JWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing Authorization header"})
		}

		// Accept "Bearer <token>" or raw token
		var tokenString string
		if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			tokenString = strings.TrimSpace(auth[7:])
		} else {
			tokenString = strings.TrimSpace(auth)
		}
		if tokenString == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing token"})
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return c.Status(500).JSON(fiber.Map{"error": "JWT secret not configured"})
		}

		// Parse token with MapClaims so we can be flexible about claim names
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// ensure method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenUnverifiable
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// --- extract user id (try multiple keys) ---
		var userID uuid.UUID
		foundUser := false
		tryKeys := []string{"user_id", "userId", "UserID", "sub"}
		for _, k := range tryKeys {
			if v, ok := claims[k]; ok && v != nil {
				if s, ok := v.(string); ok && s != "" {
					id, err := uuid.Parse(s)
					if err == nil {
						userID = id
						foundUser = true
						break
					}
				}
			}
		}
		if !foundUser {
			// sometimes numeric sub or other types — try formatting
			for _, k := range tryKeys {
				if v, ok := claims[k]; ok && v != nil {
					// try convert to string via Sprintf
					s := ""
					switch t := v.(type) {
					case string:
						s = t
					default:
						// ignore non-string if not parseable
					}
					if s != "" {
						id, err := uuid.Parse(s)
						if err == nil {
							userID = id
							foundUser = true
							break
						}
					}
				}
			}
		}
		if !foundUser {
			return c.Status(401).JSON(fiber.Map{"error": "user id not found in token"})
		}

		// --- extract role (try multiple keys) ---
		var roleName string
		var roleFound bool
		roleKeys := []string{"role", "Role"}
		for _, k := range roleKeys {
			if v, ok := claims[k]; ok && v != nil {
				if s, ok := v.(string); ok && s != "" {
					// Try if it's UUID (role id) first
					if id, err := uuid.Parse(s); err == nil {
						// role is an UUID -> lookup by ID
						roleObj, err := repository.GetRoleByID(id)
						if err == nil && roleObj != nil {
							roleName = roleObj.Name
							roleFound = true
							break
						}
					}
					// else treat s as role name -> lookup to verify it exists
					roleObj, err := repository.GetRoleByName(s)
					if err == nil && roleObj != nil {
						roleName = roleObj.Name
						roleFound = true
						break
					}
					// fallback: set roleName to claim string (even if no db match)
					roleName = s
					roleFound = true
					break
				}
			}
		}
		if !roleFound {
			// role not present -> deny
			return c.Status(401).JSON(fiber.Map{"error": "role not found in token"})
		}

		// store values in locals for handlers: user_id (uuid.UUID) and role (string)
		c.Locals("user_id", userID)
		c.Locals("role", roleName)

		return c.Next()
	}
}
