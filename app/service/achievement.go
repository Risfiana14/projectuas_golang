// app/service/achievement.go
package service

import (

	"projectuas/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"
)

// Fungsi yang sudah ada & jalan
func DeleteAchievement(c *fiber.Ctx, pgDB *sql.DB, mongoClient *mongo.Client) error {
	userID := c.Locals("user_id").(uuid.UUID)
	userRole := c.Locals("role").(string)

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "ID tidak valid", "success": false})
	}

	ref, err := repository.GetReference(pgDB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"message": "Prestasi tidak ditemukan", "success": false})
		}
		return c.Status(500).JSON(fiber.Map{"message": "Error database", "success": false})
	}

	// RBAC: hanya admin atau pemilik (mahasiswa) yang boleh hapus draft
	if userRole != "admin" {
		if userRole == "mahasiswa" && *ref.StudentID != userID {
			return c.Status(403).JSON(fiber.Map{"message": "Bukan prestasi Anda", "success": false})
		}
		if userRole == "dosen_wali" {
			return c.Status(403).JSON(fiber.Map{"message": "Dosen tidak boleh menghapus", "success": false})
		}
	}

	if *ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"message": "Hanya draft yang boleh dihapus", "success": false})
	}

	// Hapus dari MongoDB & PostgreSQL
	if err := repository.SoftDeleteMongo(mongoClient, *ref.MongoAchievementID); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal hapus MongoDB", "success": false})
	}
	if err := repository.UpdateToDeleted(pgDB, id); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal update PostgreSQL", "success": false})
	}

	return c.JSON(fiber.Map{"message": "Prestasi draft berhasil dihapus", "success": true})
}

// Fungsi dummy agar route tidak error (bisa kamu isi nanti)
func CreateAchievement(c *fiber.Ctx) error       { return c.JSON(fiber.Map{"message": "Create OK"}) }
func GetMyAchievements(c *fiber.Ctx) error         { return c.JSON(fiber.Map{"message": "My Achievements"}) }
func GetAchievementDetail(c *fiber.Ctx) error      { return c.JSON(fiber.Map{"message": "Detail"}) }
func UpdateAchievement(c *fiber.Ctx) error         { return c.JSON(fiber.Map{"message": "Update OK"}) }
func SubmitAchievement(c *fiber.Ctx) error         { return c.JSON(fiber.Map{"message": "Submitted"}) }
func GetPendingAchievements(c *fiber.Ctx) error       { return c.JSON(fiber.Map{"message": "Pending List"}) }
func VerifyAchievement(c *fiber.Ctx) error         { return c.JSON(fiber.Map{"message": "Verified"}) }
func RejectAchievement(c *fiber.Ctx) error         { return c.JSON(fiber.Map{"message": "Rejected"}) }
func GetAllAchievements(c *fiber.Ctx) error        { return c.JSON(fiber.Map{"message": "All Achievements"}) }
func GetDashboardStats(c *fiber.Ctx) error         { return c.JSON(fiber.Map{"stats": "100 prestasi"}) }