package repository

import (
    "projectuas/app/model"
    "github.com/google/uuid"
    "time"
)

func GetAchievementRefsByStudentID(studentID uuid.UUID) ([]model.AchievementRef, error) {
    rows, err := DB.Query(`
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at,
               verified_by, rejection_note, created_at, updated_at
        FROM achievement_references
        WHERE student_id=$1
        ORDER BY created_at DESC
    `, studentID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []model.AchievementRef
    for rows.Next() {
        var ref model.AchievementRef
        _ = rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoID, &ref.Status,
            &ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
            &ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt)
        list = append(list, ref)
    }
    return list, nil
}

func CreateAchievementReference(id uuid.UUID, studentID uuid.UUID, mongoID string) error {
    _, err := DB.Exec(`
        INSERT INTO achievement_references 
        (id, student_id, mongo_achievement_id, status, created_at, updated_at)
        VALUES ($1, $2, $3, 'draft', NOW(), NOW())
    `, id, studentID, mongoID)
    return err
}

func UpdateAchievementRef(ref *model.AchievementRef) error {
    _, err := DB.Exec(`
        UPDATE achievement_references 
        SET status=$1, submitted_at=$2, verified_at=$3, verified_by=$4,
            rejection_note=$5, updated_at=$6
        WHERE id=$7
    `,
        ref.Status, ref.SubmittedAt, ref.VerifiedAt, ref.VerifiedBy,
        ref.RejectionNote, time.Now(), ref.ID)
    return err
}