// app/model/achievement.go
package model

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

type Claims struct {
    UserID uuid.UUID `json:"user_id"`
    Role   string    `json:"role"`
    jwt.RegisteredClaims
}

type Achievement struct {
    ID          string    `json:"id" bson:"_id,omitempty"`
    Title       string    `json:"title" bson:"title"`
    Description string    `json:"description" bson:"description"`
    Category    string    `json:"category" bson:"category"`
    Level       string    `json:"level" bson:"level"`
    AwardDate   string    `json:"award_date" bson:"award_date"`
    Organizer   string    `json:"organizer" bson:"organizer"`
    StudentID   uuid.UUID `bson:"student_id"`
    Status      string    `json:"status" bson:"status"`
    CreatedAt   time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
    Attachments []Attachment `bson:"attachments"`
}

type AchievementRef struct {
    ID            uuid.UUID  `db:"id"`
    StudentID     uuid.UUID  `db:"student_id"`
    MongoID       string     `db:"mongo_achievement_id"`
    Status        string     `db:"status"`
    SubmittedAt   *time.Time `db:"submitted_at"`
    VerifiedAt    *time.Time `db:"verified_at"`
    VerifiedBy    *uuid.UUID `db:"verified_by"`
    RejectionNote *string    `db:"rejection_note"`
    CreatedAt     time.Time  `db:"created_at"`
    UpdatedAt     time.Time  `db:"updated_at"`
}

type AchievementHistory struct {
	ID          uuid.UUID  `json:"id"`
	ReferenceID uuid.UUID  `json:"reference_id"`
	Status      string     `json:"status"`
	Note        *string    `json:"note,omitempty"`     // ⬅️ pointer
	UserID      *uuid.UUID `json:"user_id,omitempty"` // ⬅️ pointer
	Timestamp   time.Time  `json:"timestamp"`
}

type Attachment struct {
    FileName   string    `bson:"fileName" json:"fileName"`
    FileURL    string    `bson:"fileUrl" json:"fileUrl"`
    FileType   string    `bson:"fileType" json:"fileType"`
    UploadedAt time.Time `bson:"uploadedAt" json:"uploadedAt"`
}