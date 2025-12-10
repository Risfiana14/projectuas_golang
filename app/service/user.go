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
        RoleID   string `json:"role_id"`
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
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return fiber.NewError(400, "invalid user id")
    }

    user, err := repository.GetUserByID(id)
    if err != nil {
        return fiber.NewError(404, "user not found")
    }

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

    user.Username = req.Username
    user.Email = req.Email
    user.FullName = req.FullName
    user.RoleID = uuid.MustParse(req.RoleID)
    user.IsActive = req.IsActive

    if err := repository.UpdateUser(user); err != nil {
        return fiber.NewError(500, err.Error())
    }

    return c.JSON(fiber.Map{"message": "User updated", "user": user})
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
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return fiber.NewError(400, "invalid user id")
    }

    user, err := repository.GetUserByID(id)
    if err != nil {
        return fiber.NewError(404, "user not found")
    }

    var req struct {
        RoleID string `json:"role_id"`
    }

    if err := c.BodyParser(&req); err != nil {
        return fiber.NewError(400, err.Error())
    }

    user.RoleID = uuid.MustParse(req.RoleID)

    if err := repository.UpdateUser(user); err != nil {
        return fiber.NewError(500, err.Error())
    }

    return c.JSON(fiber.Map{"message": "Role updated", "user": user})
}

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

