package repository

import (
    "database/sql"
    "projectuas/app/model"
    "github.com/google/uuid"
)

func GetReference(db *sql.DB, id uuid.UUID) (*model.AchievementReference, error) {
    ref := &model.AchievementReference{}
    err := db.QueryRow(`
        SELECT id, student_id, mongo_achievement_id, status, submitted_at,
               verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references WHERE id = $1`, id).
        Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
            &ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
            &ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt)
    return ref, err
}

func UpdateToDeleted(db *sql.DB, id uuid.UUID) error {
    _, err := db.Exec(`UPDATE achievement_references SET status='deleted', updated_at=NOW() WHERE id=$1`, id)
    return err
}