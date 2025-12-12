package repository

import (
	"projectuas/app/model"
	"github.com/google/uuid"
)

// GetLecturers returns list of lecturers
func GetLecturers() ([]model.Lecturer, error) {
	rows, err := DB.Query(`
        SELECT id, user_id, lecturer_id, department, created_at
        FROM lecturers
        ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Lecturer
	for rows.Next() {
		var l model.Lecturer
		// adjust scanning to match model.Lecturer fields: ID, UserID, LecturerID, Department, CreatedAt
		if err := rows.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, nil
}

// GetAdviseesByLecturerID returns students that have this lecturer as advisor
func GetAdviseesByLecturerID(lecturerID uuid.UUID) ([]model.Student, error) {
	rows, err := DB.Query(`
        SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
        FROM students
        WHERE advisor_id = $1
    `, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Student
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

// GetLecturerByUserID finds lecturer row by users.id (user_id)
func GetLecturerByUserID(userID uuid.UUID) (*model.Lecturer, error) {
	var l model.Lecturer
	err := DB.QueryRow(`
        SELECT id, user_id, lecturer_id, department, created_at
        FROM lecturers
        WHERE user_id = $1
    `, userID).Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &l, nil
}
