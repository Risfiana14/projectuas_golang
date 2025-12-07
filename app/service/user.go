// app/service/user.go
package service

import (
    "projectuas/app/model"
    "projectuas/app/repository"
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)

func GetUsers(c *fiber.Ctx) error {
    users, err := repository.GetAllUsers()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
    }
    return c.JSON(fiber.Map{"status": "success", "data": users})
}

func CreateUserAdmin(c *fiber.Ctx) error {
    type Req struct {
        Username string `json:"username"`
        Password string `json:"password"`
        FullName string `json:"fullName"`
        Role     string `json:"role"`
    }
    var req Req
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Format salah"})
    }

    user := &model.User{
        Username: req.Username,
        FullName: req.FullName,
        Role:     req.Role,
    }
    if err := repository.CreateUser(user, req.Password); err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal buat user"})
    }
    return c.JSON(fiber.Map{"status": "success", "message": "User berhasil dibuat", "data": user})
}

func UpdateUserRole(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, _ := uuid.Parse(idStr)
    type Req struct{ Role string `json:"role"` }
    var req Req
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Role wajib diisi"})
    }
    if err := repository.UpdateUserRole(id, req.Role); err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
    }
    return c.JSON(fiber.Map{"status": "success", "message": "Role berhasil diupdate"})
}

func DeleteUserAdmin(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, _ := uuid.Parse(idStr)
    if err := repository.DeleteUser(id); err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
    }
    return c.JSON(fiber.Map{"status": "success", "message": "User berhasil dihapus"})
}

func Profile(c *fiber.Ctx) error {
    // Ambil user_id dari JWT middleware
    userID := c.Locals("user_id").(uuid.UUID)

    user, err := repository.GetUserByID(userID.String())
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "status":  "error",
            "message": "User tidak ditemukan",
        })
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "data": fiber.Map{
            "user": fiber.Map{
                "id":        user.ID.String(),
                "username":  user.Username,
                "fullName":  user.FullName,
                "role":      user.Role,
            },
        },
    })
}