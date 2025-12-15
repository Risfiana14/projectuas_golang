package service

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"projectuas/app/repository"
	"projectuas/app/model"
)

// GET STUDENTS
func GetStudents(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	// 🔐 RBAC sesuai SRS
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "forbidden",
		})
	}

	list, err := repository.GetStudents()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed get students")
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   list,
	})
}

// GET STUDENT DETAIL
func GetStudentDetail(c *fiber.Ctx) error {
    studentID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Invalid student ID",
		})
	}

    st, err := repository.GetStudentByID(studentID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Student not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": st,
	})
}

// GET STUDENT ACHIEVEMENTS
func GetStudentAchievements(c *fiber.Ctx) error {
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(uuid.UUID)

    studentID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status": "error",
            "message": "Invalid student ID",
        })
    }

    // 🔐 RBAC
    if role == "mahasiswa" {
        return c.Status(403).JSON(fiber.Map{
            "status": "error",
            "message": "Forbidden",
        })
    }

    if role == "dosen_wali" {
        student, err := repository.GetStudentByID(studentID)
        if err != nil || student == nil {
            return c.Status(404).JSON(fiber.Map{
                "status": "error",
                "message": "Student not found",
            })
        }

        if student.AdvisorID == nil || *student.AdvisorID != userID {
            return c.Status(403).JSON(fiber.Map{
                "status": "error",
                "message": "You are not the advisor of this student",
            })
        }
    }

    refs, err := repository.GetAchievementRefsByStudentID(studentID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status": "error",
            "message": "Failed to get achievements",
        })
    }

    // ✅ PASTIKAN TIDAK NIL
    if refs == nil {
        refs = []model.AchievementRef{}
    }

    // ✅ TAMBAHAN MESSAGE JIKA KOSONG
    if len(refs) == 0 {
        return c.JSON(fiber.Map{
            "status":  "success",
            "message": "No achievements available with current filter",
            "data":    refs,
        })
    }

    // ✅ NORMAL RESPONSE
    return c.JSON(fiber.Map{
        "status": "success",
        "data":   refs,
    })
}

// ASSIGN ADVISOR
func AssignAdvisor(c *fiber.Ctx) error {
    studentID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status": "error",
            "message": "invalid student id",
        })
    }

    var body struct {
        AdvisorID string `json:"advisor_id"`
    }

    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status": "error",
            "message": "invalid request body",
        })
    }

    advisorUUID, err := uuid.Parse(body.AdvisorID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status": "error",
            "message": "invalid advisor id",
        })
    }

    if err := repository.AssignAdvisorToStudent(studentID, advisorUUID); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status": "error",
            "message": err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "message": "advisor assigned successfully",
    })
}
