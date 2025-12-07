// app/model/user.go
package model

import "github.com/google/uuid"

type User struct {
    ID           uuid.UUID `json:"id" db:"id"`
    Username     string    `json:"username" db:"username"`
    PasswordHash string    `json:"-" db:"password_hash"`  // PASTI db:"password_hash"
    FullName     string    `json:"fullName" db:"full_name"`
    Role         string    `json:"role" db:"role"`
}