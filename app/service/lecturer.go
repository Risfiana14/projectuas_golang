package service

import (
    "net/http"

    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"

    "projectuas/app/repository"
    "projectuas/app/model"
)

// GET LECTURERS
// GetLecturers godoc
// @Summary Get lecturers
// @Description Get list of lecturers
// @Tags Lecturers
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.Lecturer
// @Router /lecturers [get]
func GetLecturers(c *fiber.Ctx) error {
    list, err := repository.GetLecturers()
    if err != nil {
        return fiber.NewError(http.StatusInternalServerError, "failed get lecturers")
    }
    return c.JSON(list)
}

// GET ADVISEES BY LECTURER
// GetLecturerAdvisees godoc
// @Summary Get lecturer advisees
// @Description Lecturer view list of advised students
// @Tags Lecturers
// @Security BearerAuth
// @Param id path string true "Lecturer ID"
// @Produce json
// @Success 200 {array} model.Student
// @Failure 403 {object} map[string]string
// @Router /lecturers/{id}/advisees [get]
func GetLecturerAdvisees(c *fiber.Ctx) error {
    lectID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return fiber.NewError(http.StatusBadRequest, "invalid id")
    }

    role := c.Locals("role").(string)
    callerID := c.Locals("user_id").(uuid.UUID)

    if role == "dosen_wali" {
        if callerID != lectID {
            return fiber.NewError(http.StatusForbidden, "forbidden: cannot view advisees of other lecturers")
        }
    }

    list, err := repository.GetAdviseesByLecturerID(lectID)
    if err != nil {
        return fiber.NewError(http.StatusInternalServerError, "failed get advisees")
    }

    return c.JSON(list)
}

// GET PENDING ACHIEVEMENTS

// GetPendingAchievements godoc
// @Summary Get pending achievements
// @Description Admin sees all submitted achievements, dosen wali sees only advisees' submitted achievements
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.AchievementRef
// @Failure 403 {object} map[string]string
// @Router /achievements/pending [get]
func GetPendingAchievements(c *fiber.Ctx) error {
    role := c.Locals("role").(string)
    callerID := c.Locals("user_id").(uuid.UUID)

    refs, err := repository.GetAchievementsByStatus("submitted")
    if err != nil {
        return fiber.NewError(http.StatusInternalServerError, "failed fetch achievements")
    }

    // Admin -> langsung return semua
    if role == "admin" {
        return c.JSON(refs)
    }

    // Dosen wali -> filter prestasi mahasiswa bimbingannya
    if role == "dosen_wali" {
        filtered := []*model.AchievementRef{} // memakai pointer

        for _, ref := range refs {
            student, err := repository.GetStudentByUserID(ref.StudentID)
            if err != nil || student == nil {
                continue
            }

            if student.AdvisorID != nil && *student.AdvisorID == callerID {
                filtered = append(filtered, ref)
            }
        }

        return c.JSON(filtered)
    }

    return fiber.NewError(http.StatusForbidden, "insufficient role")
}


// FR-006: View Prestasi Mahasiswa Bimbingan

// GetLecturerAchievements godoc
// @Summary View achievements of a lecturer's advisees
// @Description Lecturer can view submitted achievements of their advisees
// @Tags Achievements
// @Security BearerAuth
// @Param id path string true "Lecturer ID"
// @Produce json
// @Success 200 {array} model.AchievementRef
// @Failure 403 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /lecturers/{id}/achievements [get]
func GetLecturerAchievements(c *fiber.Ctx) error {
    lectID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return fiber.NewError(400, "invalid lecturer id")
    }

    role := c.Locals("role").(string)
    callerID := c.Locals("user_id").(uuid.UUID)

    // Dosen wali hanya bisa melihat prestasi mahasiswa bimbingannya sendiri
    if role == "dosen_wali" && callerID != lectID {
        return fiber.NewError(403, "forbidden: cannot view other lecturer's data")
    }

    // Get list mahasiswa bimbingan
    advisees, err := repository.GetAdviseesByLecturerID(lectID)
    if err != nil {
        return fiber.NewError(500, "failed to get advisees")
    }

    // Ambil semua student_id
    studentIDs := []uuid.UUID{}
    for _, s := range advisees {
        studentIDs = append(studentIDs, s.ID)
    }

    // Ambil prestasi milik mahasiswa tersebut
    achievements, err := repository.GetAchievementsByStudents(studentIDs)
    if err != nil {
        return fiber.NewError(500, "failed to get achievements")
    }

    return c.JSON(achievements)
}