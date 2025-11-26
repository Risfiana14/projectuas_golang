package model

import (
    "time"
    "github.com/google/uuid"
    "github.com/golang-jwt/jwt/v5"
)

type AchievementReference struct {
    ID                 *uuid.UUID `json:"id"`
    StudentID          *uuid.UUID `json:"student_id"`
    MongoAchievementID *string    `json:"mongo_achievement_id"`
    Status             *string    `json:"status"` // draft, submitted, verified, rejected, deleted
    SubmittedAt        *time.Time `json:"submitted_at"`
    VerifiedAt         *time.Time `json:"verified_at"`
    VerifiedBy         *uuid.UUID `json:"verified_by"`
    RejectionNote      *string    `json:"rejection_note"`
    CreatedAt          *time.Time `json:"created_at"`
    UpdatedAt          *time.Time `json:"updated_at"`
}

type Claims struct {
    UserID uuid.UUID `json:"user_id"`
    Role   string    `json:"role"` // admin, mahasiswa, dosen_wali
    jwt.RegisteredClaims
}