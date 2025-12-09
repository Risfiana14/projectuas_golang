package service

import (
    "projectuas/app/model"
    "projectuas/app/repository"
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "log"
)

// GET ALL USERS
func GetUsers(c *fiber.Ctx) error {
    users, err := repository.GetAllUsers()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
    }
    return c.JSON(fiber.Map{"status": "success", "data": users})
}

// CREATE USER
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

    role, err := repository.GetRoleByName(req.Role)
    if err != nil {
        log.Println("ERROR GET ROLE:", err)
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal ambil role"})
    }
    if role == nil {
        return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Role tidak valid"})
    }

    user := &model.User{
        Username: req.Username,
        FullName: req.FullName,
        RoleID:   role.ID,
    }

    if err := repository.CreateUser(user, req.Password); err != nil {
        log.Println("ERROR CREATE USER:", err)
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal buat user"})
    }

    return c.JSON(fiber.Map{"status": "success", "message": "User berhasil dibuat", "data": user})
}

// UPDATE USER ROLE
func UpdateUserRole(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"status": "error", "message": "ID tidak valid"})
    }

    type Req struct{ Role string `json:"role"` }
    var req Req
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Role wajib diisi"})
    }

    role, err := repository.GetRoleByName(req.Role)
    if err != nil || role == nil {
        return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Role tidak valid"})
    }

    if err := repository.UpdateUserRole(id, role.ID); err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
    }

    return c.JSON(fiber.Map{"status": "success", "message": "Role berhasil diupdate"})
}

// DELETE USER
func DeleteUserAdmin(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"status": "error", "message": "ID tidak valid"})
    }

    if err := repository.DeleteUser(id); err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
    }

    return c.JSON(fiber.Map{"status": "success", "message": "User berhasil dihapus"})
}

// GET USER PROFILE
func Profile(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uuid.UUID)

    user, err := repository.GetUserByID(userID)
    if err != nil || user == nil {
        return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User tidak ditemukan"})
    }

    role, _ := repository.GetRoleByID(user.RoleID)
    var roleName string
    if role != nil {
        roleName = role.Name
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "data": fiber.Map{
            "user": fiber.Map{
                "id":       user.ID.String(),
                "username": user.Username,
                "fullName": user.FullName,
                "role":     roleName,
            },
        },
    })
}
