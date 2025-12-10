package model

import (
    "time"
    "github.com/google/uuid"
)

type Lecturer struct {
    ID        uuid.UUID `db:"id" json:"id"`
    UserID    uuid.UUID `db:"user_id" json:"user_id"`
    FullName  string    `db:"full_name" json:"full_name"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}
