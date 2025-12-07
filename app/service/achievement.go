// app/service/achievement.go
package service

import (
    "context"
    "os"
    "time"
    "projectuas/app/model"
    "projectuas/app/repository"
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateAchievement(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uuid.UUID)
    var ach model.Achievement
    if err := c.BodyParser(&ach); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
    }

    ach.StudentID = userID
    ach.Status = "draft"
    ach.CreatedAt = time.Now()
    ach.UpdatedAt = time.Now()

    coll := repository.MongoClient.Database(os.Getenv("DB_NAME")).Collection("achievements")
    res, err := coll.InsertOne(context.TODO(), ach)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal simpan MongoDB"})
    }

    mongoID := res.InsertedID.(primitive.ObjectID).Hex()
    refID := uuid.New()

    _, err = repository.DB.Exec(`
        INSERT INTO achievement_references (id, student_id, mongo_achievement_id, status)
        VALUES ($1, $2, $3, 'draft')`, refID, userID, mongoID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal simpan reference"})
    }

    return c.JSON(fiber.Map{
        "message":  "Prestasi draft berhasil dibuat",
        "ref_id":   refID,
        "mongo_id": mongoID,
    })
}

// Dummy handler (biar route tidak error)
func GetMyAchievements(c *fiber.Ctx) error      { return c.JSON(fiber.Map{"message": "Daftar prestasi saya"}) }
func GetAchievementDetail(c *fiber.Ctx) error   { return c.JSON(fiber.Map{"message": "Detail"}) }
func UpdateAchievement(c *fiber.Ctx) error      { return c.JSON(fiber.Map{"message": "Updated"}) }
func DeleteAchievement(c *fiber.Ctx) error      { return c.JSON(fiber.Map{"message": "Deleted"}) }
func SubmitAchievement(c *fiber.Ctx) error      { return c.JSON(fiber.Map{"message": "Submitted"}) }
func GetPendingAchievements(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "Pending list"}) }
func VerifyAchievement(c *fiber.Ctx) error      { return c.JSON(fiber.Map{"message": "Verified"}) }
func RejectAchievement(c *fiber.Ctx) error      { return c.JSON(fiber.Map{"message": "Rejected"}) }
func GetAllAchievements(c *fiber.Ctx) error     { return c.JSON(fiber.Map{"message": "All achievements"}) }