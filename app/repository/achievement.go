package repository

import (
	"errors"
	"github.com/google/uuid"
	"projectuas/app/model"
	"time"
)

// GET REF BY ID
func GetAchievementRefByID(id uuid.UUID) (*model.AchievementRef, error) {
	row := DB.QueryRow(`
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at,
		       verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE id=$1
	`, id)

	var ref model.AchievementRef
	err := row.Scan(&ref.ID, &ref.StudentID, &ref.MongoID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
		&ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

// GET REF BY USER (alias GetAchievementRefsByUser)
func GetAchievementRefsByUser(userID uuid.UUID) ([]*model.AchievementRef, error) {
	refs, err := GetAchievementRefsByStudentID(userID)
	if err != nil {
		return nil, err
	}
	var ptrs []*model.AchievementRef
	for i := range refs {
		ptrs = append(ptrs, &refs[i])
	}
	return ptrs, nil
}

// GET ALL REFS
func GetAllAchievementRefs() ([]*model.AchievementRef, error) {
	rows, err := DB.Query(`
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at,
		       verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.AchievementRef
	for rows.Next() {
		var ref model.AchievementRef
		_ = rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoID, &ref.Status,
			&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
			&ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt)
		list = append(list, &ref)
	}
	return list, nil
}

// GET BY STATUS
func GetAchievementsByStatus(status string) ([]*model.AchievementRef, error) {
	rows, err := DB.Query(`
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at,
		       verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE status=$1
		ORDER BY created_at DESC
	`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.AchievementRef
	for rows.Next() {
		var ref model.AchievementRef
		_ = rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoID, &ref.Status,
			&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
			&ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt)
		list = append(list, &ref)
	}
	return list, nil
}

// DELETE REFERENCE
func DeleteAchievementRef(id uuid.UUID) error {
	res, err := DB.Exec(`DELETE FROM achievement_references WHERE id=$1`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("reference not found")
	}
	return nil
}

// ADD HISTORY
func AddAchievementHistory(mongoID string, action string, userID uuid.UUID) error {
	_, err := DB.Exec(`
		INSERT INTO achievement_history (id, mongo_achievement_id, action, user_id, timestamp)
		VALUES ($1, $2, $3, $4, $5)
	`, uuid.New(), mongoID, action, userID, time.Now())
	return err
}

// SAVE ATTACHMENT
func SaveAttachment(mongoID string, filename string, data []byte) (string, error) {
	// simpan file ke folder lokal atau cloud storage
	// disini kita buat dummy URL
	url := "https://example.com/files/" + filename
	return url, nil
}


func GetAchievementHistory(mongoID string) ([]model.AchievementHistory, error) {
    rows, err := DB.Query(`
        SELECT id, mongo_achievement_id, action, user_id, timestamp
        FROM achievement_history
        WHERE mongo_achievement_id=$1
        ORDER BY timestamp ASC
    `, mongoID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var history []model.AchievementHistory
    for rows.Next() {
        var h model.AchievementHistory
        _ = rows.Scan(&h.ID, &h.MongoID, &h.Action, &h.UserID, &h.Timestamp)
        history = append(history, h)
    }
    return history, nil
}

func GetAchievementStatistics() (map[string]int, error) {
    rows, err := DB.Query(`
        SELECT status, COUNT(*) 
        FROM achievement_references
        GROUP BY status
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    stats := make(map[string]int)
    for rows.Next() {
        var status string
        var count int
        _ = rows.Scan(&status, &count)
        stats[status] = count
    }
    return stats, nil
}

func GetStudentReport(studentID uuid.UUID) ([]model.AchievementRef, error) {
    return GetAchievementRefsByStudentID(studentID)
}
