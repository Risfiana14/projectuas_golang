package service

import (
    "github.com/gofiber/fiber/v2"
    "golang.org/x/crypto/bcrypt"
    "projectuas/app/model"
    "projectuas/app/repository"
    "github.com/google/uuid"
)

// ================= PASSWORD UTILS =================
func hashPassword(pw string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(pw), 12)
    return string(b), err
}

// ===================================================
// =============== CRUD USER (ADMIN ONLY) ============
// ===================================================

// GET /api/v1/users
func AdminGetUsers(c *fiber.Ctx) error {
    users, err := repository.GetAllUsers()
    if err != nil {
        return fiber.NewError(500, err.Error())
    }
    return c.JSON(users)
}

// GET /api/v1/users/:id
func AdminGetUserDetail(c *fiber.Ctx) error {
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return fiber.NewError(400, "invalid user id")
    }

    user, err := repository.GetUserByID(id)
    if err != nil {
        return fiber.NewError(404, "user not found")
    }

    return c.JSON(user)
}

// POST /api/v1/users
func AdminCreateUser(c *fiber.Ctx) error {
    var req struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
        FullName string `json:"full_name"`
        RoleID   string `json:"roleId"`
        IsActive bool   `json:"is_active"`
    }

    if err := c.BodyParser(&req); err != nil {
        return fiber.NewError(400, err.Error())
    }

    hashed, err := hashPassword(req.Password)
    if err != nil {
        return fiber.NewError(500, "failed hash password")
    }

    user := &model.User{
        ID:           uuid.New(),
        Username:     req.Username,
        Email:        req.Email,
        FullName:     req.FullName,
        PasswordHash: hashed,
        RoleID:       uuid.MustParse(req.RoleID),
        IsActive:     req.IsActive,
    }

    if err := repository.CreateUser(user); err != nil {
        return fiber.NewError(500, err.Error())
    }

    return c.JSON(fiber.Map{"message": "User created", "user": user})
}

// PUT /api/v1/users/:id
func AdminUpdateUser(c *fiber.Ctx) error {
    // Parse userID dari URL
    userID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return fiber.NewError(400, "invalid user id")
    }

    // Ambil user dari DB
    user, err := repository.GetUserByID(userID)
    if err != nil {
        return fiber.NewError(404, "user not found")
    }

    // Parse request body
    var req struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        FullName string `json:"full_name"`
        RoleID   string `json:"role_id"`
        IsActive bool   `json:"is_active"`
    }

    if err := c.BodyParser(&req); err != nil {
        return fiber.NewError(400, err.Error())
    }

    // Update fields user jika ada
    if req.Username != "" {
        user.Username = req.Username
    }
    if req.Email != "" {
        user.Email = req.Email
    }
    if req.FullName != "" {
        user.FullName = req.FullName
    }
    if req.RoleID != "" {
        roleUUID, err := uuid.Parse(req.RoleID)
        if err != nil {
            return fiber.NewError(400, "invalid role id")
        }
        user.RoleID = roleUUID
    }
    user.IsActive = req.IsActive

    // Update ke database
    if err := repository.UpdateUser(user); err != nil {
        return fiber.NewError(500, err.Error())
    }

    return c.JSON(fiber.Map{
        "message": "User updated",
        "user": fiber.Map{
            "id":        user.ID.String(),
            "username":  user.Username,
            "email":     user.Email,
            "full_name": user.FullName,
            "role_id":   user.RoleID.String(),
            "is_active": user.IsActive,
        },
    })
}

// DELETE /api/v1/users/:id
func AdminDeleteUser(c *fiber.Ctx) error {
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return fiber.NewError(400, "invalid user id")
    }

    if err := repository.DeleteUser(id); err != nil {
        return fiber.NewError(404, "user not found")
    }

    return c.JSON(fiber.Map{"message": "User deleted"})
}

// PUT /api/v1/users/:id/role
func AdminUpdateUserRole(c *fiber.Ctx) error {
    // Parse user ID dari URL
    userIDParam := c.Params("id")
    userID, err := uuid.Parse(userIDParam)
    if err != nil {
        return fiber.NewError(400, "invalid user id")
    }

    // Ambil user dari database
    user, err := repository.GetUserByID(userID)
    if err != nil {
        return fiber.NewError(404, "user not found")
    }

    // Parse body request
    var req struct {
        RoleID string `json:"role_id"`
    }

    if err := c.BodyParser(&req); err != nil {
        return fiber.NewError(400, "invalid request body")
    }

    // Validasi RoleID
    if req.RoleID == "" {
        return fiber.NewError(400, "role_id is required")
    }

    roleUUID, err := uuid.Parse(req.RoleID)
    if err != nil {
        return fiber.NewError(400, "invalid role_id")
    }

    // Update role user
    user.RoleID = roleUUID

    if err := repository.UpdateUser(user); err != nil {
        return fiber.NewError(500, "failed to update user role")
    }

    return c.JSON(fiber.Map{
        "message": "User role updated successfully",
        "user": fiber.Map{
            "id":     user.ID.String(),
            "username": user.Username,
            "roleId": user.RoleID.String(),
        },
    })
}

// GET /api/v1/profile
func Profile(c *fiber.Ctx) error {
    userID := c.Locals("user_id")

    id, ok := userID.(uuid.UUID)
    if !ok {
        return c.Status(400).JSON(fiber.Map{
            "error": "invalid user id",
        })
    }

    user, err := repository.UserFindByID(id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "error": "user not found",
        })
    }

    return c.JSON(fiber.Map{
        "id":        user.ID,
        "username":  user.Username,
        "email":     user.Email,
        "full_name": user.FullName,
        "role_id":   user.RoleID,
        "is_active": user.IsActive,
    })
}

