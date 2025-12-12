package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"projectuas/app/model"
)

// NOTE: file ini mengasumsikan ada var DB *sql.DB dideklarasikan di paket repository (inisialisasi di database.ConnectPostgres()).

// GET ACHIEVEMENT REFS BY STUDENT ID
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
		err := rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoID, &ref.Status,
			&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
			&ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt)
		if err != nil {
			return nil, err
		}
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
		if err == sql.ErrNoRows {
			return nil, err
		}
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
		err := rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoID, &ref.Status,
			&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
			&ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt)
		if err != nil {
			return nil, err
		}
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
		ORDER BY submitted_at DESC
	`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.AchievementRef
	for rows.Next() {
		var ref model.AchievementRef
		err := rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoID, &ref.Status,
			&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
			&ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt)
		if err != nil {
			return nil, err
		}
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
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("reference not found")
	}
	return nil
}

// AddAchievementHistory inserts a history row (Postgres)
func AddAchievementHistory(refID uuid.UUID, status string, note string, userID uuid.UUID) error {
    var userPtr *uuid.UUID
    if userID != uuid.Nil {
        userPtr = &userID
    } else {
        userPtr = nil // supaya masuk NULL, tidak FK error
    }

    _, err := DB.Exec(`
        INSERT INTO achievement_history (reference_id, status, note, user_id, timestamp)
        VALUES ($1, $2, $3, $4, NOW())
    `, refID, status, note, userPtr)

    return err
}

// SAVE ATTACHMENT (dummy)
func SaveAttachment(mongoID string, filename string, data []byte) (string, error) {
	// dummy URL
	url := "https://example.com/files/" + filename
	return url, nil
}

// GET ACHIEVEMENT HISTORY (query only)
func GetAchievementHistory(refID uuid.UUID) ([]model.AchievementHistory, error) {
    rows, err := DB.Query(`
        SELECT id, reference_id, status, note, user_id, timestamp
        FROM achievement_history
        WHERE reference_id = $1
        ORDER BY timestamp ASC
    `, refID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []model.AchievementHistory
    for rows.Next() {
        var h model.AchievementHistory
        if err := rows.Scan(&h.ID, &h.ReferenceID, &h.Status, &h.Note, &h.UserID, &h.Timestamp); err != nil {
            return nil, err
        }
        list = append(list, h)
    }
    return list, nil
}

// Get statistics (simple)
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
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats[status] = count
	}
	return stats, nil
}

// GET STUDENT REPORT (returns refs)
func GetStudentReport(studentID uuid.UUID) ([]model.AchievementRef, error) {
	return GetAchievementRefsByStudentID(studentID)
}

// Get achievements by student list
func GetAchievementsByStudents(studentIDs []uuid.UUID) ([]*model.AchievementRef, error) {
	if len(studentIDs) == 0 {
		return []*model.AchievementRef{}, nil
	}

	query := `
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at,
               verified_by, rejection_note, created_at, updated_at
        FROM achievement_references
        WHERE student_id = ANY($1)
        ORDER BY created_at DESC
    `

	rows, err := DB.Query(query, pq.Array(studentIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.AchievementRef
	for rows.Next() {
		var r model.AchievementRef
		if err := rows.Scan(&r.ID, &r.StudentID, &r.MongoID, &r.Status,
			&r.SubmittedAt, &r.VerifiedAt, &r.VerifiedBy,
			&r.RejectionNote, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, &r)
	}
	return results, nil
}

// Update status helper (Postgres)
func UpdateAchievementStatus(
    refID uuid.UUID,
    status string,
    t *time.Time,
    actorID *uuid.UUID,
    note *string,
) error {

    if status == "verified" {
        _, err := DB.Exec(`
            UPDATE achievement_references
            SET status=$1, verified_at=$2, verified_by=$3, rejection_note=NULL, updated_at=NOW()
            WHERE id=$4
        `, status, t, actorID, refID)
        return err
    }

    if status == "rejected" {
        _, err := DB.Exec(`
            UPDATE achievement_references
            SET status=$1, rejection_note=$2, verified_at=NULL, verified_by=NULL, updated_at=NOW()
            WHERE id=$3
        `, status, note, refID)
        return err
    }

    return errors.New("invalid status")
}

