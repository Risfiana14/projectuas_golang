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
// CreateAchievement godoc
// @Summary Create achievement (draft)
// @Description Mahasiswa membuat prestasi dengan status draft
// @Tags Achievements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Router /achievements [post]
// @Param body body model.Achievement true "Achievement payload"
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

	// ✅ HISTORY WAJIB SESUAI SRS
	note := "Achievement created as draft"
	_ = repository.AddAchievementHistory(
		refID,
		"draft",
		&note,
		userID,
	)

	return c.Status(201).JSON(fiber.Map{
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
// SubmitAchievement godoc
// @Summary Submit achievement
// @Description Mahasiswa mengirim prestasi dari draft ke submitted
// @Tags Achievements
// @Security BearerAuth
// @Param id path string true "Achievement Reference ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /achievements/{id}/submit [post]
func SubmitAchievement(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "mahasiswa" {
		return fiber.NewError(403, "only mahasiswa can submit achievement")
	}

	userID := c.Locals("user_id").(uuid.UUID)
	refID := uuid.MustParse(c.Params("id"))

	ref, err := repository.GetAchievementRefByID(refID)
	if err != nil {
		return fiber.NewError(404, "achievement not found")
	}

	// 🔒 VALIDASI KEPEMILIKAN (WAJIB SRS)
	student, err := repository.GetStudentByUserID(userID)
	if err != nil || ref.StudentID != student.ID {
		return fiber.NewError(403, "not your achievement")
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

	// Sync MongoDB
	if err := repository.UpdateAchievementMongo(
		ref.MongoID,
		bson.M{
			"status":     "submitted",
			"updated_at": time.Now(),
		},
	); err != nil {
		return fiber.NewError(500, "failed sync mongo status")
	}

	note := "Achievement submitted by student"
	_ = repository.AddAchievementHistory(ref.ID, "submitted", &note, userID)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "achievement submitted",
	})
}

// FR-007: Verify Achievement
// VerifyAchievement godoc
// @Summary Verify achievement
// @Description Dosen wali memverifikasi prestasi mahasiswa bimbingan
// @Tags Achievements
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /achievements/{id}/verify [post]
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

	// UPDATE POSTGRES
	if err := repository.UpdateAchievementStatus(
		achievement.ID,
		"verified",
		&now,
		&userID,
		nil,
	); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Failed to update achievement status",
		})
	}

	// 🔥 WAJIB INSERT HISTORY
		note := "Achievement verified by advisor"
	if err := repository.AddAchievementHistory(
		achievement.ID,
		"verified",
		&note, // ✅ pointer ke note yang benar
		userID,
	); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Failed to insert history",
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
// RejectAchievement godoc
// @Summary Reject achievement
// @Description Dosen wali menolak prestasi mahasiswa bimbingan
// @Tags Achievements
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Accept json
// @Produce json
// @Param body body object{note=string} true "Rejection note"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /achievements/{id}/reject [post]
func RejectAchievement(c *fiber.Ctx) error {
	achID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Invalid achievement id",
		})
	}

	userID := c.Locals("user_id").(uuid.UUID)

	// ===== BODY =====
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

	// ===== GET ACHIEVEMENT =====
	achievement, err := repository.GetAchievementRefByID(achID)
	if err != nil || achievement == nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "error",
			"message": "Achievement not found",
		})
	}

	// ===== FINAL STATE GUARD (SRS) =====
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

	if achievement.Status != "submitted" {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Only submitted achievements can be rejected",
		})
	}

	// ===== VALIDASI DOSEN WALI =====
	student, err := repository.GetStudentByID(achievement.StudentID)
	if err != nil || student == nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "error",
			"message": "Student not found",
		})
	}

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

	// ===== UPDATE STATUS DI POSTGRES =====
	if err := repository.UpdateAchievementStatus(
		achievement.ID,
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

	// 🔥 WAJIB INSERT HISTORY
	if err := repository.AddAchievementHistory(
		achievement.ID,
		"rejected",
		&body.Note,
		userID,
	); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Failed to insert history",
		})
	}

	// ===== SYNC KE MONGODB =====
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
// GetAchievementDetail godoc
// @Summary Get achievement detail
// @Description Detail prestasi (mahasiswa, dosen wali, admin)
// @Tags Achievements
// @Security BearerAuth
// @Param id path string true "Achievement Reference ID"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Router /achievements/{id} [get]
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
		return fiber.NewError(500, "failed to fetch achievement")
	}

	if ach.Attachments == nil {
		ach.Attachments = []model.Attachment{}
	}

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(uuid.UUID)

	switch role {
	case "mahasiswa":
		student, err := repository.GetStudentByUserID(userID)
		if err != nil || ref.StudentID != student.ID {
			return c.Status(403).JSON(fiber.Map{
				"status":  "error",
				"message": "Not your achievement",
			})
		}

	case "dosen_wali":
		student, err := repository.GetStudentByID(ref.StudentID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Student not found",
			})
		}
		lecturer, err := repository.GetLecturerByUserID(userID)
		if err != nil || student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
			return c.Status(403).JSON(fiber.Map{
				"status":  "error",
				"message": "Not your advisee's achievement",
			})
		}

	case "admin":
		// allowed
	default:
		return c.Status(403).JSON(fiber.Map{
			"status":  "error",
			"message": "Forbidden",
		})
	}

	return c.JSON(fiber.Map{"reference": ref, "achievement": ach})
}


// GET ALL (filter by role)
// GetAllAchievements godoc
// @Summary Get all achievements
// @Description Admin melihat semua, dosen wali melihat mahasiswa bimbingan
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.AchievementRef
// @Failure 403 {object} map[string]string
// @Router /achievements [get]
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
		refs, err = repository.GetAchievementsByAdvisorUserID(userID)
	
	case "mahasiswa":
		// 🔥 INI YANG HILANG
		student, err := repository.GetStudentByUserID(userID)
		if err != nil {
			return c.Status(403).JSON(fiber.Map{
				"status":  "error",
				"message": "student not found",
			})
		}

		refs, err = repository.GetAchievementRefsByStudentID(student.ID)
		
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

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   refs,
	})
}

// UPDATE (only draft)
// UpdateAchievement godoc
// @Summary Update achievement
// @Description Mahasiswa mengedit prestasi (hanya draft)
// @Tags Achievements
// @Security BearerAuth
// @Param id path string true "Achievement Reference ID"
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /achievements/{id} [put]
// @Param body body model.Achievement true "Updated achievement data"
func UpdateAchievement(c *fiber.Ctx) error {
	refID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(400, "invalid reference id")
	}

	ref, err := repository.GetAchievementRefByID(refID)
	if err != nil {
		return fiber.NewError(404, "Reference not found")
	}

	if ref.Status != "draft" {
		return fiber.NewError(400, "Hanya draft yang bisa diedit")
	}

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(uuid.UUID)

	if role != "mahasiswa" {
		return fiber.NewError(403, "only mahasiswa can update achievement")
	}

	student, err := repository.GetStudentByUserID(userID)
	if err != nil || ref.StudentID != student.ID {
		return fiber.NewError(403, "not your achievement")
	}

	var body model.Achievement
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(400, "Invalid JSON")
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
		return fiber.NewError(500, "Gagal update MongoDB")
	}

	return c.JSON(fiber.Map{"message": "Berhasil update prestasi"})
}

// DELETE (only draft)
// DeleteAchievement godoc
// @Summary Delete achievement
// @Description Hapus prestasi (soft delete, hanya draft)
// @Tags Achievements
// @Security BearerAuth
// @Param id path string true "Achievement Reference ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /achievements/{id} [delete]
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

	// ===============================
	// 🔒 RULE 1: STATUS (GLOBAL)
	// ===============================
	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Only achievements with status 'draft' can be deleted",
		})
	}

	// ===============================
	// 🔐 RULE 2: ROLE
	// ===============================
	switch role {

	case "mahasiswa":
		student, err := repository.GetStudentByUserID(userID)
		if err != nil {
			return fiber.NewError(403, "student record not found")
		}

		if ref.StudentID != student.ID {
			return c.Status(403).JSON(fiber.Map{
				"status":  "error",
				"message": "cannot delete other student's achievement",
			})
		}

	case "admin":
		// ✅ admin boleh, TAPI tetap hanya draft

	default:
		// ❌ dosen_wali & role lain
		return fiber.NewError(403, "forbidden")
	}

	// ===============================
	// 🗑️ SOFT DELETE
	// ===============================

	// 1️⃣ Soft delete di MongoDB
	if err := repository.SoftDeleteAchievementMongo(ref.MongoID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to soft delete achievement in MongoDB",
		})
	}

	// 2️⃣ Update reference di PostgreSQL
	if err := repository.SoftDeleteAchievementRef(ref.ID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to update achievement reference",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Achievement deleted successfully",
	})
}

// GET ACHIEVEMENT HISTORY (PG)
// GetAchievementHistory godoc
// @Summary Get achievement history
// @Description Riwayat perubahan status prestasi
// @Tags Achievements
// @Security BearerAuth
// @Param id path string true "Achievement Reference ID"
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Router /achievements/{id}/history [get]
func GetAchievementHistory(c *fiber.Ctx) error {
	refID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid reference id",
		})
	}

	ref, err := repository.GetAchievementRefByID(refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "reference not found",
		})
	}

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(uuid.UUID)

	// 🔐 RBAC
	switch role {

	case "mahasiswa":
		student, err := repository.GetStudentByUserID(userID)
		if err != nil || ref.StudentID != student.ID {
			return c.Status(403).JSON(fiber.Map{
				"status":  "error",
				"message": "not your achievement",
			})
		}

	case "dosen_wali":
		student, err := repository.GetStudentByID(ref.StudentID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "student not found",
			})
		}

		lecturer, err := repository.GetLecturerByUserID(userID)
		if err != nil || student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
			return c.Status(403).JSON(fiber.Map{
				"status":  "error",
				"message": "not your advisee",
			})
		}

	case "admin":
		// allowed

	default:
		return c.Status(403).JSON(fiber.Map{
			"status":  "error",
			"message": "forbidden",
		})
	}

	history, err := repository.GetAchievementHistoryByRefID(ref.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to fetch history",
		})
	}

	// ✅ RESPONSE SESUAI SRS
	if len(history) == 0 {
		return c.JSON(fiber.Map{
			"status": "success",
			"data":   []interface{}{},
		})
	}

	// Mapping response agar rapi
	var response []fiber.Map
	for _, h := range history {
		response = append(response, fiber.Map{
			"id":        h.ID,
			"status":    h.Status,
			"note":      h.Note,
			"actor_id":  h.UserID,
			"timestamp": h.Timestamp,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

// UPLOAD ATTACHMENTS
// UploadAchievementAttachments godoc
// @Summary Upload achievement attachments
// @Description Upload file pendukung prestasi (hanya status draft)
// @Tags Achievements
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Param files formData file true "Attachment files (multiple)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /achievements/{id}/attachments [post]
func UploadAchievementAttachments(c *fiber.Ctx) error {
	refID := c.Params("id")
	if refID == "" {
		return fiber.NewError(400, "invalid achievement id")
	}

	userID := c.Locals("user_id").(uuid.UUID)
	role := c.Locals("role").(string)

	// 🔒 ROLE HARUS MAHASISWA
	if role != "mahasiswa" {
		return fiber.NewError(403, "only mahasiswa can upload attachments")
	}

	ref, err := repository.GetAchievementRefByID(uuid.MustParse(refID))
	if err != nil {
		return fiber.NewError(404, "achievement not found")
	}

	// 🔒 STATUS HARUS DRAFT
	if ref.Status != "draft" {
		return fiber.NewError(400, "attachments can only be uploaded for draft achievement")
	}

	// 🔒 VALIDASI KEPEMILIKAN
	student, err := repository.GetStudentByUserID(userID)
	if err != nil || ref.StudentID != student.ID {
		return fiber.NewError(403, "not your achievement")
	}

	form, err := c.MultipartForm()
	if err != nil || len(form.File["files"]) == 0 {
		return fiber.NewError(400, "no files uploaded")
	}

	ach, err := repository.GetAchievementMongoByID(ref.MongoID)
	if err != nil {
		return fiber.NewError(500, "failed to get achievement")
	}

	if ach.Attachments == nil {
		ach.Attachments = []model.Attachment{}
	}

	for _, fileHeader := range form.File["files"] {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}

		func() {
			defer file.Close()
			data := make([]byte, fileHeader.Size)
			file.Read(data)

			url, err := repository.SaveAttachment(ref.MongoID, fileHeader.Filename, data)
			if err != nil {
				return
			}

			ach.Attachments = append(ach.Attachments, model.Attachment{
				FileName:   fileHeader.Filename,
				FileURL:    url,
				FileType:   fileHeader.Header.Get("Content-Type"),
				UploadedAt: time.Now(),
			})
		}()
	}

	if err := repository.UpdateAchievementMongo(ref.MongoID, bson.M{
		"attachments": ach.Attachments,
		"updated_at":  time.Now(),
	}); err != nil {
		return fiber.NewError(500, "failed to save attachments")
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "attachments uploaded successfully",
		"data":    ach.Attachments,
	})
}
