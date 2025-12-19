package service

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"projectuas/app/repository"
)

// GET ACHIEVEMENT STATISTICS
// GetAchievementStatistics godoc
// @Summary Achievement statistics
// @Description Get achievement analytics
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /reports/statistics [get]
func GetAchievementStatistics(c *fiber.Ctx) error {
	stats, err := repository.GetAchievementStatistics()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed get stats")
	}
	return c.JSON(stats)
}

// GET STUDENT REPORT

// GetStudentReport godoc
// @Summary Get report of a student
// @Description Get submitted achievements of a specific student
// @Tags Reports
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reports/students/{id} [get]
func GetStudentReport(c *fiber.Ctx) error {
	studentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "invalid student id",
		})
	}

	refs, err := repository.GetStudentReport(studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "failed to get student report",
		})
	}

	if len(refs) == 0 {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "No achievements available for this student",
			"data":    []interface{}{},
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   refs,
	})
}

