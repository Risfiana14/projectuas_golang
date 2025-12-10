package service

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"projectuas/app/repository"
)

// GET LECTURERS
func GetLecturers(c *fiber.Ctx) error {
	list, err := repository.GetLecturers()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed get lecturers")
	}
	return c.JSON(list)
}

// GET ADVISEES BY LECTURER
func GetLecturerAdvisees(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid id")
	}
	list, err := repository.GetAdviseesByLecturerID(id)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed get advisees")
	}
	return c.JSON(list)
}

// GET PENDING ACHIEVEMENTS (for lecturers)
func GetPendingAchievements(c *fiber.Ctx) error {
	refs, err := repository.GetAchievementsByStatus("submitted")
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Gagal mengambil prestasi pending")
	}
	return c.JSON(refs)
}
