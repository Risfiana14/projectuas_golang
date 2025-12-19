package service

import (
	"log"
	"os"
	"projectuas/app/model"
	"projectuas/app/repository"
	"time"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// --- LOGIN REQUEST STRUCT (global, untuk Swagger) ---
type LoginRequest struct {
	Username string `json:"username" example:"user123"`
	Password string `json:"password" example:"secret"`
}

// Login godoc
// @Summary Login user
// @Description Login dan mendapatkan JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println("ERROR: Body parser:", err)
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Format salah"})
	}

	log.Println("LOGIN REQUEST:", req.Username)

	user, err := repository.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		log.Println("USER TIDAK DITEMUKAN:", req.Username)
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Username atau password salah"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Println("PASSWORD SALAH:", err)
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Username atau password salah"})
	}

	role, err := repository.GetRoleByID(user.RoleID)
	if err != nil {
		log.Println("ERROR GET ROLE:", err)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Internal server error"})
	}

	permissions, err := repository.GetPermissionsByRoleID(user.RoleID)
	if err != nil {
		log.Println("ERROR GET PERMISSIONS:", err)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Internal server error"})
	}

	// --- Create JWT ---
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Claims{
		UserID: user.ID,
		Role:   role.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	// --- Create Refresh Token ---
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Claims{
		UserID: user.ID,
		Role:   role.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	})
	refreshTokenString, _ := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token":        tokenString,
			"refreshToken": refreshTokenString,
			"user": fiber.Map{
				"id":          user.ID.String(),
				"username":    user.Username,
				"fullName":    user.FullName,
				"role":        role.Name,
				"permissions": permissions,
			},
		},
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout user and invalidate token (optional blacklist)
// @Tags Authentication
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/logout [post]
func Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Missing Authorization header"})
	}
	_ = strings.TrimPrefix(authHeader, "Bearer ")
	// optional: add token blacklist logic here
	return c.JSON(fiber.Map{"status": "success", "message": "Logout berhasil"})
}

// Refresh godoc
// @Summary Refresh JWT token
// @Description Refresh access token using current valid token
// @Tags Authentication
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func Refresh(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Token invalid"})
	}

	claims := token.Claims.(*model.Claims)

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Claims{
		UserID: claims.UserID,
		Role:   claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	newTokenString, _ := newToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token": newTokenString,
		},
	})
}
