package service

import (
	"io"
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
	var ach model.Achievement
	if err := c.BodyParser(&ach); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid JSON")
	}

	ach.StudentID = userID
	ach.Status = "draft"
	ach.CreatedAt = time.Now()
	ach.UpdatedAt = time.Now()

	mongoID, err := repository.InsertAchievementMongo(&ach)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Gagal simpan MongoDB")
	}

	refID := uuid.New()
	if err := repository.CreateAchievementReference(refID, userID, mongoID); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Gagal simpan reference")
	}

	return c.JSON(fiber.Map{"message": "Draft prestasi berhasil dibuat", "ref_id": refID, "mongo_id": mongoID})
}

// SUBMIT
func SubmitAchievement(c *fiber.Ctx) error {
	refID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid reference id")
	}

	ref, err := repository.GetAchievementRefByID(refID)
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "Reference not found")
	}

	if ref.Status != "draft" {
		return fiber.NewError(http.StatusBadRequest, "Only draft can be submitted")
	}

	now := time.Now()
	ref.Status = "submitted"
	ref.SubmittedAt = &now
	if err := repository.UpdateAchievementRef(ref); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Gagal update status")
	}

	_ = repository.AddAchievementHistory(ref.MongoID, "submitted", ref.StudentID)
	return c.JSON(fiber.Map{"message": "Prestasi berhasil di-submit"})
}

// VERIFY
func VerifyAchievement(c *fiber.Ctx) error {
	refID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid reference id")
	}

	ref, err := repository.GetAchievementRefByID(refID)
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "Reference not found")
	}

	if ref.Status != "submitted" {
		return fiber.NewError(http.StatusBadRequest, "Only submitted can be verified")
	}

	userID := c.Locals("user_id").(uuid.UUID)
	now := time.Now()
	ref.Status = "verified"
	ref.VerifiedBy = &userID
	ref.VerifiedAt = &now

	if err := repository.UpdateAchievementRef(ref); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Gagal verifikasi")
	}

	_ = repository.AddAchievementHistory(ref.MongoID, "verified", userID)
	return c.JSON(fiber.Map{"message": "Prestasi berhasil diverifikasi"})
}

// REJECT
func RejectAchievement(c *fiber.Ctx) error {
	refID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid reference id")
	}

	var body struct {
		Note string `json:"note"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid JSON")
	}

	ref, err := repository.GetAchievementRefByID(refID)
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "Reference not found")
	}

	now := time.Now()
	ref.Status = "rejected"
	ref.RejectionNote = &body.Note
	ref.VerifiedAt = &now

	if err := repository.UpdateAchievementRef(ref); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Gagal menolak prestasi")
	}

	_ = repository.AddAchievementHistory(ref.MongoID, "rejected", uuid.Nil)
	return c.JSON(fiber.Map{"message": "Prestasi ditolak"})
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
	refs, err := repository.GetAchievementRefsByUser(userID)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Gagal mengambil data")
	}
	return c.JSON(refs)
}

// GET ALL (filter by role)
func GetAllAchievements(c *fiber.Ctx) error {
	roleI := c.Locals("role")
	role, _ := roleI.(string)
	userID, _ := c.Locals("user_id").(uuid.UUID)

	switch role {
	case "mahasiswa":
		refs, err := repository.GetAchievementRefsByUser(userID)
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, "Gagal mengambil data")
		}
		return c.JSON(refs)
	case "dosen_wali":
		refs, err := repository.GetAchievementsByStatus("submitted")
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, "Gagal mengambil data")
		}
		return c.JSON(refs)
	default: // admin
		refs, err := repository.GetAllAchievementRefs()
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, "Gagal mengambil semua prestasi")
		}
		return c.JSON(refs)
	}
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
		return fiber.NewError(http.StatusBadRequest, "invalid reference id")
	}

	ref, err := repository.GetAchievementRefByID(refID)
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "Reference not found")
	}

	if ref.Status != "draft" {
		return fiber.NewError(http.StatusBadRequest, "Hanya draft yang bisa dihapus")
	}

	if err := repository.DeleteAchievementMongo(ref.MongoID); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Gagal hapus MongoDB")
	}

	if err := repository.DeleteAchievementRef(refID); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Gagal hapus reference")
	}

	return c.JSON(fiber.Map{"message": "Berhasil menghapus prestasi"})
}

// GET HISTORY
func GetAchievementHistory(c *fiber.Ctx) error {
	refID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid reference id")
	}
	ref, err := repository.GetAchievementRefByID(refID)
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "Reference not found")
	}
	history, err := repository.GetAchievementHistory(ref.MongoID)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Gagal ambil history")
	}
	return c.JSON(history)
}

// UPLOAD ATTACHMENTS
func UploadAchievementAttachments(c *fiber.Ctx) error {
	refID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid reference id")
	}
	ref, err := repository.GetAchievementRefByID(refID)
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "Reference not found")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "file required")
	}
	f, err := file.Open()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed open file")
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed read file")
	}

	url, err := repository.SaveAttachment(ref.MongoID, file.Filename, data)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed save attachment")
	}
	_ = repository.AddAchievementHistory(ref.MongoID, "attachment_uploaded", uuid.Nil)

	return c.JSON(fiber.Map{"message": "Uploaded", "url": url})
}
