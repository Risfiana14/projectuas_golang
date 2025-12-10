package service

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"projectuas/app/repository"
)

// GET ACHIEVEMENT STATISTICS
func GetAchievementStatistics(c *fiber.Ctx) error {
	stats, err := repository.GetAchievementStatistics()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed get stats")
	}
	return c.JSON(stats)
}

// GET STUDENT REPORT
func GetStudentReport(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid id")
	}
	report, err := repository.GetStudentReport(id)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed get report")
	}
	return c.JSON(report)
}
