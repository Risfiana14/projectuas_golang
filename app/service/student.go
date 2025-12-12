package service

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"projectuas/app/repository"
)

// GET STUDENTS
func GetStudents(c *fiber.Ctx) error {
	list, err := repository.GetStudents()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed get students")
	}
	return c.JSON(list)
}

// GET STUDENT DETAIL
func GetStudentDetail(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid id")
	}
	st, err := repository.GetStudentByUserID(id)
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "student not found")
	}
	return c.JSON(st)
}

// GET STUDENT ACHIEVEMENTS
func GetStudentAchievements(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid id")
	}
	refs, err := repository.GetAchievementRefsByStudentID(id)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed get achievements")
	}
	return c.JSON(refs)
}

// ASSIGN ADVISOR
func AssignAdvisor(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid id")
	}
	var body struct {
		AdvisorID string `json:"advisor_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid body")
	}
	advisorID, err := uuid.Parse(body.AdvisorID)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid advisor id")
	}
	if err := repository.AssignAdvisorToStudent(id, advisorID); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed assign advisor")
	}
	return c.JSON(fiber.Map{"message": "Advisor assigned"})
}
