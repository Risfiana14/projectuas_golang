package repository

import (
    "projectuas/app/model"
    "github.com/google/uuid"
    "time"
)

// GetStudents
func GetStudents() ([]model.Student, error) {
    rows, err := DB.Query(`
        SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
        FROM students ORDER BY created_at DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []model.Student
    for rows.Next() {
        var s model.Student
        _ = rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt)
        list = append(list, s)
    }
    return list, nil
}

// GET STUDENT BY ID
func GetStudentByID(studentID uuid.UUID) (*model.Student, error) {
	row := DB.QueryRow(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id
		FROM students
		WHERE id = $1
	`, studentID)

	var student model.Student
	err := row.Scan(
		&student.ID,
		&student.UserID,
        &student.StudentID,
		&student.ProgramStudy,
		&student.AcademicYear,
		&student.AdvisorID,
	)

	if err != nil {
		return nil, err
	}

	return &student, nil
}
// GET STUDENT BY USER ID
func GetStudentByUserID(userID uuid.UUID) (*model.Student, error) {
    var s model.Student
    err := DB.QueryRow(`
        SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
        FROM students 
        WHERE user_id = $1
    `, userID).Scan(
        &s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy,
        &s.AcademicYear, &s.AdvisorID, &s.CreatedAt,
    )

    if err != nil {
        return nil, err
    }
    return &s, nil
}


func AssignAdvisorToStudent(studentID uuid.UUID, advisorID uuid.UUID) error {
    _, err := DB.Exec(`UPDATE students SET advisor_id=$1, updated_at=$2 WHERE id=$3`,
        advisorID, time.Now(), studentID)
    return err
}
