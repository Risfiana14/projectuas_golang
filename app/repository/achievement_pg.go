// app/repository/achievement_pg.go
package repository

import (
    "context"
    "database/sql"
    "os"
    "time"

    "projectuas/app/model"
    "github.com/google/uuid"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

var DB *sql.DB
var MongoClient *mongo.Client

// Init dipanggil dari route/setup.go
func Init(db *sql.DB, client *mongo.Client) {
    DB = db
    MongoClient = client
}

// --- PostgreSQL Functions ---
func GetRef(id uuid.UUID) (*model.AchievementRef, error) {
    ref := &model.AchievementRef{}
    err := DB.QueryRow(`
        SELECT id, student_id, mongo_achievement_id, status, submitted_at,
               verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references WHERE id = $1`, id).
        Scan(&ref.ID, &ref.StudentID, &ref.MongoID, &ref.Status,
            &ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
            &ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt)
    return ref, err
}

func SubmitAchievementPG(id uuid.UUID) error {
    _, err := DB.Exec(`UPDATE achievement_references SET status='submitted', submitted_at=NOW(), updated_at=NOW() WHERE id=$1`, id)
    return err
}

func UpdateToDeleted(id uuid.UUID) error {
    _, err := DB.Exec(`UPDATE achievement_references SET status='deleted', updated_at=NOW() WHERE id=$1`, id)
    return err
}

// --- MongoDB Functions ---
func SoftDeleteMongo(mongoID string) error {
    coll := MongoClient.Database(os.Getenv("DB_NAME")).Collection("achievements")
    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    _, err := coll.UpdateOne(ctx, bson.M{"_id": mongoID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
    return err
}