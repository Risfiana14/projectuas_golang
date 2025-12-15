package service

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"projectuas/app/model"
	"projectuas/app/repository"
)

// CREATE DRAFT
func CreateAchievement(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	role := c.Locals("role").(string)
	if role != "mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"status": "error",
			"message": "Only mahasiswa can create achievements",
		})
	}

	student, err := repository.GetStudentByUserID(userID)
	if err != nil {
		return fiber.NewError(403, "student record not found")
	}

	var ach model.Achievement
	if err := c.BodyParser(&ach); err != nil {
		return fiber.NewError(400, "invalid json")
	}

	ach.StudentID = student.ID
	ach.Status = "draft"
	ach.CreatedAt = time.Now()
	ach.UpdatedAt = time.Now()

	mongoID, err := repository.InsertAchievementMongo(&ach)
	if err != nil {
		return fiber.NewError(500, "failed save mongo")
	}

	refID := uuid.New()
	if err := repository.CreateAchievementReference(refID, student.ID, mongoID); err != nil {
		return fiber.NewError(500, "failed save reference")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"message": "Achievement draft created successfully",
		"data": fiber.Map{
			"reference_id": refID,
			"mongo_id":     mongoID,
			"status":       "draft",
			"created_at":   ach.CreatedAt,
		},
	})
}

// SUBMIT
func SubmitAchievement(c *fiber.Ctx) error {
	refID := uuid.MustParse(c.Params("id"))

	ref, err := repository.GetAchievementRefByID(refID)
	if err != nil {
		return fiber.NewError(404, "achievement not found")
	}

	if ref.Status != "draft" {
		return fiber.NewError(400, "only draft can be submitted")
	}

	now := time.Now()
	ref.Status = "submitted"
	ref.SubmittedAt = &now

	if err := repository.UpdateAchievementRef(ref); err != nil {
		return fiber.NewError(500, "failed update status")
	}
	// 🔴 UPDATE STATUS DI MONGODB
	if err := repository.UpdateAchievementMongo(
		ref.MongoID,
		bson.M{
			"status":     "submitted",
			"updated_at": time.Now(),
		},
	); err != nil {
		return fiber.NewError(500, "failed sync mongo status")
	}

	_ = repository.AddAchievementHistory(ref.ID, "submitted", "", ref.StudentID)
	return c.JSON(fiber.Map{"message": "submitted"})
}

// FR-007: Verify Achievement
func VerifyAchievement(c *fiber.Ctx) error {
	achID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Invalid achievement id",
		})
	}

	userID := c.Locals("user_id").(uuid.UUID)

	achievement, err := repository.GetAchievementRefByID(achID)
	if err != nil || achievement == nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "error",
			"message": "Achievement not found",
		})
	}

	if achievement.Status != "submitted" {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Only submitted achievements can be verified",
		})
	}

	student, err := repository.GetStudentByID(achievement.StudentID)
	if err != nil || student == nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "error",
			"message": "Student not found",
		})
	}

	// 🔧 FIX UTAMA DI SINI
	lecturer, err := repository.GetLecturerByUserID(userID)
	if err != nil || lecturer == nil {
		return c.Status(403).JSON(fiber.Map{
			"status": "error",
			"message": "Lecturer not found",
		})
	}

	if student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
		return c.Status(403).JSON(fiber.Map{
			"status": "error",
			"message": "You are not the advisor of this student",
		})
	}

	now := time.Now()
	if err := repository.UpdateAchievementStatus(
	achID,
	"verified",
	&now,
	&userID,
	nil,
); err != nil {
	return c.Status(500).JSON(fiber.Map{
		"status": "error",
		"message": "Failed to verify achievement",
	})
}

	// 🔴 UPDATE STATUS DI MONGODB
	if err := repository.UpdateAchievementMongo(
		achievement.MongoID,
		bson.M{
			"status":     "verified",
			"updated_at": time.Now(),
		},
	); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Failed to sync achievement status to MongoDB",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"message": "Achievement verified successfully",
	})
}

// FR-008: Reject Achievement
func RejectAchievement(c *fiber.Ctx) error {
    achID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status": "error",
            "message": "Invalid achievement id",
        })
    }

    userID := c.Locals("user_id").(uuid.UUID)

    type Request struct {
        Note string `json:"note"`
    }

    var body Request
    if err := c.BodyParser(&body); err != nil || body.Note == "" {
        return c.Status(422).JSON(fiber.Map{
            "status": "error",
            "message": "Rejection note is required",
        })
    }

    achievement, err := repository.GetAchievementRefByID(achID)
    if err != nil || achievement == nil {
        return c.Status(404).JSON(fiber.Map{
            "status": "error",
            "message": "Achievement not found",
        })
    }

    // GUARD FINAL STATE (SRS)
	if achievement.Status == "verified" {
		return c.Status(409).JSON(fiber.Map{
			"status": "error",
			"message": "Verified achievement cannot be rejected",
		})
	}

	if achievement.Status == "rejected" {
		return c.Status(409).JSON(fiber.Map{
			"status": "error",
			"message": "Achievement already rejected",
		})
	}

	// Only submitted allowed
	if achievement.Status != "submitted" {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Only submitted achievements can be rejected",
		})
	}

    student, err := repository.GetStudentByID(achievement.StudentID)
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

    now := time.Now()

    if err := repository.UpdateAchievementStatus(
		achID,
		"rejected",
		&now,
		&userID,
		&body.Note,
	); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Failed to reject achievement",
		})
	}

	// 🔴 UPDATE STATUS DI MONGODB
	if err := repository.UpdateAchievementMongo(
		achievement.MongoID,
		bson.M{
			"status":     "rejected",
			"updated_at": time.Now(),
		},
	); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Failed to sync achievement status to MongoDB",
		})
	}

    return c.JSON(fiber.Map{
        "status": "success",
        "message": "Achievement rejected successfully",
    })
}

// GET DETAIL
func GetAchievementDetail(c *fiber.Ctx) error {
	refID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid reference id")
	}

	ref, err := repository.GetAchievementRefByID(refID)
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "Reference not found")
	}

	ach, err := repository.GetAchievementMongoByID(ref.MongoID)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Gagal mengambil data MongoDB")
	}

	return c.JSON(fiber.Map{"reference": ref, "achievement": ach})
}

// GET MY ACHIEVEMENTS
func GetMyAchievements(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	student, err := repository.GetStudentByUserID(userID)
	if err != nil {
		return fiber.NewError(403, "student not found")
	}

	refs, _ := repository.GetAchievementRefsByStudentID(student.ID)
	return c.JSON(refs)
}

// GET ALL (filter by role)
func GetAllAchievements(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(uuid.UUID)

	var (
		refs []*model.AchievementRef
		err  error
	)

	switch role {

	case "admin":
		refs, err = repository.GetAllAchievementRefs()

	case "dosen_wali":
		// 1. Ambil mahasiswa bimbingan dosen
		students, err := repository.GetStudentsByAdvisor(userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": "error",
				"message": "failed to get advisees",
			})
		}

		var studentIDs []uuid.UUID
		for _, s := range students {
			studentIDs = append(studentIDs, s.ID)
		}

		// 2. Ambil achievement mereka
		refs, err = repository.GetAchievementsByStudents(studentIDs)

	default:
		return c.Status(403).JSON(fiber.Map{
			"status": "error",
			"message": "forbidden",
		})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "failed to get achievements",
		})
	}

	if len(refs) == 0 {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "No achievements available",
			"data":    []interface{}{},
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   refs,
	})
}


// UPDATE (only draft)
func UpdateAchievement(c *fiber.Ctx) error {
	refID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid reference id")
	}

	ref, err := repository.GetAchievementRefByID(refID)
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "Reference not found")
	}

	if ref.Status != "draft" {
		return fiber.NewError(http.StatusBadRequest, "Hanya draft yang bisa diedit")
	}

	var body model.Achievement
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid JSON")
	}

	body.UpdatedAt = time.Now()
	update := bson.M{
		"title":       body.Title,
		"description": body.Description,
		"category":    body.Category,
		"level":       body.Level,
		"award_date":  body.AwardDate,
		"organizer":   body.Organizer,
		"updated_at":  body.UpdatedAt,
	}

	if err := repository.UpdateAchievementMongo(ref.MongoID, update); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Gagal update MongoDB")
	}

	return c.JSON(fiber.Map{"message": "Berhasil update prestasi"})
}

// DELETE (only draft)
func DeleteAchievement(c *fiber.Ctx) error {
	refID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(400, "invalid reference id")
	}

	userID := c.Locals("user_id").(uuid.UUID)
	role := c.Locals("role").(string)

	ref, err := repository.GetAchievementRefByID(refID)
	if err != nil {
		return fiber.NewError(404, "achievement not found")
	}

	if role == "mahasiswa" {
		student, err := repository.GetStudentByUserID(userID)
		if err != nil {
			return fiber.NewError(403, "student record not found")
		}

		if ref.StudentID != student.ID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "cannot delete other student's achievement",
			})
		}

		if ref.Status != "draft" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Only achievements with status 'draft' can be deleted",
			})
		}
	}

	if role != "admin" && role != "mahasiswa" {
		return fiber.NewError(403, "forbidden")
	}

	_ = repository.DeleteAchievementMongo(ref.MongoID)
	_ = repository.DeleteAchievementHistoryByRefID(ref.ID)
	_ = repository.DeleteAchievementRef(ref.ID)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "achievement deleted successfully",
	})
}


// GET ACHIEVEMENT HISTORY (PG)
func GetAchievementHistory(c *fiber.Ctx) error {
	refID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid reference id",
		})
	}

	ref, err := repository.GetAchievementRefByID(refID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "reference not found",
		})
	}

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(uuid.UUID)

	switch role {

	case "mahasiswa":
		student, err := repository.GetStudentByUserID(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student record not found",
			})
		}

		if ref.StudentID != student.ID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "not your achievement",
			})
		}

	case "dosen_wali":
		student, err := repository.GetStudentByID(ref.StudentID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student not found",
			})
		}
		lecturer, err := repository.GetLecturerByUserID(userID)
		if err != nil {
			return fiber.NewError(403, "lecturer record not found")
		}
		if student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
			return fiber.NewError(403, "not your advisee")
		}

	case "admin":
		// allowed

	default:
		return fiber.NewError(403, "role not allowed")
	}

	history, err := repository.GetAchievementHistory(ref.ID)
	if err != nil {
		return fiber.NewError(500, "failed to fetch history")
	}

	return c.JSON(history)
}

// UPLOAD ATTACHMENTS
func UploadAchievementAttachments(c *fiber.Ctx) error {
	refID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(400, "invalid achievement id")
	}

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(uuid.UUID)

	ref, err := repository.GetAchievementRefByID(refID)
	if err != nil {
		return fiber.NewError(404, "achievement not found")
	}

	if role == "mahasiswa" {
		student, err := repository.GetStudentByUserID(userID)
		if err != nil {
			return fiber.NewError(403, "student record not found")
		}

		if ref.StudentID != student.ID {
			return c.Status(403).JSON(fiber.Map{
				"status":  "error",
				"message": "cannot upload attachment to other student's achievement",
			})
		}
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "attachment uploaded successfully",
	})
}

