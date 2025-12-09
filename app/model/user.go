package model

import "github.com/google/uuid"

type User struct {
    ID           uuid.UUID `json:"id" db:"id"`
    Username     string    `json:"username" db:"username"`
    PasswordHash string    `json:"-" db:"password_hash"`
    FullName     string    `json:"fullName" db:"full_name"`
    RoleID       uuid.UUID `json:"roleID" db:"role_id"`
}

type Role struct {
    ID   uuid.UUID `json:"id" db:"id"`
    Name string    `json:"name" db:"name"`
}

type UserWithRole struct {
    User
    RoleName string `json:"role"`  // ambil dari tabel roles.name
}
