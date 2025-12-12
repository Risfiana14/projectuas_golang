package model

import (
    "time"
    "github.com/google/uuid"
)

type Lecturer struct {
    ID         uuid.UUID `db:"id"`
    UserID     uuid.UUID `db:"user_id"`
    LecturerID string    `db:"lecturer_id"`
    Department string    `db:"department"`
    CreatedAt  time.Time `db:"created_at"`
}
