package service

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"projectuas/app/model"
	"projectuas/app/repository"
	"github.com/google/uuid"
	"net/http"
)

func hashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), 12)
	return string(b), err
}

// GET /api/v1/users
// AdminGetUsers godoc
// @Summary Get all users
// @Description Admin get list of users
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.User
// @Router /users [get]
func AdminGetUsers(c *fiber.Ctx) error {
	users, err := repository.GetAllUsers()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(users)
}

// GET /api/v1/users/:id

// AdminGetUserDetail godoc
// @Summary Get user detail
// @Description Admin get detail of a specific user
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func AdminGetUserDetail(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid user id")
	}

	user, err := repository.GetUserByID(id)
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "user not found")
	}

	return c.JSON(user)
}

// POST /api/v1/users

// AdminCreateUser godoc
// @Summary Create new user
// @Description Admin creates a new user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body object{username=string,email=string,password=string,full_name=string,role_id=string,is_active=boolean} true "User data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func AdminCreateUser(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
		RoleID   string `json:"role_id"` // pastikan sesuai dengan JSON input
		IsActive bool   `json:"is_active"`
	}

	// Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid JSON: "+err.Error())
	}

	// Validasi input wajib
	if req.Username == "" || req.Email == "" || req.Password == "" || req.RoleID == "" {
		return fiber.NewError(http.StatusBadRequest, "username, email, password, and role_id are required")
	}

	// Hash password
	hashed, err := hashPassword(req.Password)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to hash password")
	}

	// Parse RoleID UUID
	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid role_id")
	}

	// Buat user baru
	user := &model.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		FullName:     req.FullName,
		PasswordHash: hashed,
		RoleID:       roleID,
		IsActive:     req.IsActive,
	}

	// Simpan ke repository / database
	if err := repository.CreateUser(user); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to create user: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "User created successfully",
		"user":    user,
	})
}


// PUT /api/v1/users/:id

// AdminUpdateUser godoc
// @Summary Update user
// @Description Admin updates an existing user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body object{username=string,email=string,full_name=string,role_id=string,is_active=boolean} true "Updated user data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [put]
func AdminUpdateUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid user id")
	}

	user, err := repository.GetUserByID(userID)
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "user not found")
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		FullName string `json:"full_name"`
		RoleID   string `json:"role_id"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

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
			return fiber.NewError(http.StatusBadRequest, "invalid role id")
		}
		user.RoleID = roleUUID
	}
	user.IsActive = req.IsActive

	if err := repository.UpdateUser(user); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
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

// AdminDeleteUser godoc
// @Summary Delete user
// @Description Admin deletes a user by ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
func AdminDeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid user id")
	}

	if err := repository.DeleteUser(id); err != nil {
		return fiber.NewError(http.StatusNotFound, "user not found")
	}

	return c.JSON(fiber.Map{"message": "User deleted"})
}

// PUT /api/v1/users/:id/role

// AdminUpdateUserRole godoc
// @Summary Update user role
// @Description Admin updates role of a user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param body body object{role_id=string} true "New role ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/role [put]
func AdminUpdateUserRole(c *fiber.Ctx) error {
	userIDParam := c.Params("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid user id")
	}

	user, err := repository.GetUserByID(userID)
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "user not found")
	}

	var req struct {
		RoleID string `json:"role_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid request body")
	}

	if req.RoleID == "" {
		return fiber.NewError(http.StatusBadRequest, "role_id is required")
	}

	roleUUID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid role_id")
	}

	user.RoleID = roleUUID

	if err := repository.UpdateUser(user); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed to update user role")
	}

	return c.JSON(fiber.Map{
		"message": "User role updated successfully",
		"user": fiber.Map{
			"id":       user.ID.String(),
			"username": user.Username,
			"roleId":   user.RoleID.String(),
		},
	})
}

// GET /api/v1/auth/profile

// Profile godoc
// @Summary Get profile
// @Description Get logged-in user profile
// @Tags Authentication
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /auth/profile [get]
func Profile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	id, ok := userID.(uuid.UUID)
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	user, err := repository.UserFindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
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
