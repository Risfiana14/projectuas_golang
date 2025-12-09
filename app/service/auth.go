// app/service/auth.go
package service

import (
    "log"
    "os"
    "time"
    "projectuas/app/model"
    "projectuas/app/repository"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
    type LoginRequest struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    var req LoginRequest
    if err := c.BodyParser(&req); err != nil {
        log.Println("ERROR: Body parser:", err)
        return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Format salah"})
    }

    log.Println("LOGIN REQUEST:", req.Username)

    user, err := repository.GetUserByUsername(req.Username)
    if err != nil {
        log.Println("ERROR GET USER:", err)
        return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Username atau password salah"})
    }

    if user == nil {
        log.Println("USER TIDAK DITEMUKAN:", req.Username)
        return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Username atau password salah"})
    }

    log.Println("USER DITEMUKAN:", user.Username)
    log.Println("HASH DARI DB:", user.PasswordHash)

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        log.Println("PASSWORD SALAH:", err)
        return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Username atau password salah"})
    }

    log.Println("LOGIN BERHASIL:", req.Username)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Claims{
        UserID: user.ID,
        Role:   user.RoleID.String(),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    })

    tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

    return c.JSON(fiber.Map{
        "status": "success",
        "data": fiber.Map{
            "token": tokenString,
            "user": fiber.Map{
                "id":       user.ID.String(),
                "username": user.Username,
                "fullName": user.FullName,
                "role":     user.RoleID,
            },
        },
    })
}