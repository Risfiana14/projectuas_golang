// app/service/auth.go
package service

import (
    "log" // TAMBAHKAN INI UNTUK DEBUG
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
        return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Format salah"})
    }

    // INI YANG KRUSIAL — CEK user != nil DAN err
   user, err := repository.GetUserByUsername(req.Username)
    if err != nil {
    return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error"})
    }
    if user == nil {
        return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Username atau password salah"})
    }

    log.Println("USER DITEMUKAN:", user.Username) // DEBUG

    // Bandingkan password
    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
    if err != nil {
        log.Println("PASSWORD SALAH UNTUK:", req.Username) // DEBUG
        return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Username atau password salah"})
    }

    log.Println("LOGIN BERHASIL:", req.Username) // DEBUG

    // Buat token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Claims{
        UserID: user.ID,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    })

    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal buat token"})
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "data": fiber.Map{
            "token": tokenString,
            "user": fiber.Map{
                "id":       user.ID.String(),
                "username": user.Username,
                "fullName": user.FullName,
                "role":     user.Role,
            },
        },
    })
}