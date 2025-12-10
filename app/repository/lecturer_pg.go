package repository

import (
    "projectuas/app/model"
    "github.com/google/uuid"
)

func GetLecturers() ([]model.Lecturer, error) {
    rows, err := DB.Query(`
        SELECT id, user_id, full_name, created_at 
        FROM lecturers ORDER BY created_at DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []model.Lecturer
    for rows.Next() {
        var l model.Lecturer
        _ = rows.Scan(&l.ID, &l.UserID, &l.FullName, &l.CreatedAt)
        list = append(list, l)
    }
    return list, nil
}

func GetAdviseesByLecturerID(lecturerID uuid.UUID) ([]model.Student, error) {
    rows, err := DB.Query(`
        SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at 
        FROM students WHERE advisor_id=$1`, lecturerID)
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
