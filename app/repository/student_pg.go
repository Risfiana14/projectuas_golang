package repository

import (
    "projectuas/app/model"
    "github.com/google/uuid"
    "errors"
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

// // GET STUDENTS BY ADVISOR
// func GetStudentsByAdvisor(advisorUserID uuid.UUID) ([]model.Student, error) {
// 	rows, err := DB.Query(`
// 		SELECT s.id, s.user_id, s.student_id, s.program_study,
// 		       s.academic_year, s.advisor_id, s.created_at
// 		FROM students s
// 		JOIN lecturers l ON s.advisor_id = l.id
// 		WHERE l.user_id = $1
// 	`, advisorUserID)

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var list []model.Student
// 	for rows.Next() {
// 		var s model.Student
// 		if err := rows.Scan(
// 			&s.ID, &s.UserID, &s.StudentID,
// 			&s.ProgramStudy, &s.AcademicYear,
// 			&s.AdvisorID, &s.CreatedAt,
// 		); err != nil {
// 			return nil, err
// 		}
// 		list = append(list, s)
// 	}
// 	return list, nil
// }

// ASSIGN ADVISOR TO STUDENT
func AssignAdvisorToStudent(studentID uuid.UUID, advisorID uuid.UUID) error {
    var exists bool

    err := DB.QueryRow(`
        SELECT EXISTS (
            SELECT 1 FROM lecturers WHERE id = $1
        )
    `, advisorID).Scan(&exists)

    if err != nil {
        return err
    }

    if !exists {
        return errors.New("advisor not found")
    }

    res, err := DB.Exec(`
        UPDATE students
        SET advisor_id = $1
        WHERE id = $2
    `, advisorID, studentID)

    if err != nil {
        return err
    }

    rows, _ := res.RowsAffected()
    if rows == 0 {
        return errors.New("student not found")
    }

    return nil
}

// GET STUDENTS BY ADVISOR ID
func GetStudentsByAdvisorID(advisorID uuid.UUID) ([]model.Student, error) {
    rows, err := DB.Query(`
        SELECT id, user_id, advisor_id
        FROM students
        WHERE advisor_id = $1
    `, advisorID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var students []model.Student
    for rows.Next() {
        var s model.Student
        if err := rows.Scan(&s.ID, &s.UserID, &s.AdvisorID); err != nil {
            return nil, err
        }
        students = append(students, s)
    }

    return students, nil
}
