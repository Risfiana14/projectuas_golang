package repository

import (
    "database/sql"
    "errors"
    "projectuas/app/model"
    "github.com/google/uuid"
)

var DB *sql.DB

// ============ INIT DB ============
func Init(pg *sql.DB) {
    DB = pg
}

// ============ CREATE USER ============
func CreateUser(user *model.User) error {
    query := `
        INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    _, err := DB.Exec(query,
        user.ID,
        user.Username,
        user.Email,
        user.PasswordHash,
        user.FullName,
        user.RoleID,
        user.IsActive,
    )
    return err
}

// ============ GET ALL USERS ============
func GetAllUsers() ([]model.User, error) {
    query := `SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at FROM users`

    rows, err := DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []model.User
    for rows.Next() {
        var u model.User
        err := rows.Scan(
            &u.ID,
            &u.Username,
            &u.Email,
            &u.FullName,
            &u.RoleID,
            &u.IsActive,
            &u.CreatedAt,
            &u.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    return users, nil
}

// ============ GET USER BY ID ============
func GetUserByID(id uuid.UUID) (*model.User, error) {
    query := `
        SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
        FROM users
        WHERE id = $1
    `
    var u model.User
    err := DB.QueryRow(query, id).Scan(
        &u.ID,
        &u.Username,
        &u.Email,
        &u.PasswordHash,
        &u.FullName,
        &u.RoleID,
        &u.IsActive,
        &u.CreatedAt,
        &u.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &u, nil
}

func GetUserByUsername(username string) (*model.User, error) {
    var u model.User
    err := DB.QueryRow(`
        SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
        FROM users WHERE username=$1
    `, username).Scan(
        &u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName,
        &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
    )
    return &u, err
}

func UserFindByID(id uuid.UUID) (*model.User, error) {
    query := `
        SELECT id, username, email, password_hash, full_name, role_id,
               is_active, created_at, updated_at
        FROM users
        WHERE id = $1
    `

    var u model.User
    err := DB.QueryRow(query, id).Scan(
        &u.ID,
        &u.Username,
        &u.Email,
        &u.PasswordHash,
        &u.FullName,
        &u.RoleID,
        &u.IsActive,
        &u.CreatedAt,
        &u.UpdatedAt,
    )

    if err != nil {
        return nil, errors.New("user not found")
    }

    return &u, nil
}

// ============ UPDATE USER ============
func UpdateUser(user *model.User) error {
    query := `
        UPDATE users
        SET username=$1, email=$2, full_name=$3, role_id=$4, is_active=$5, updated_at=NOW()
        WHERE id=$6
    `
    _, err := DB.Exec(query,
        user.Username,
        user.Email,
        user.FullName,
        user.RoleID,
        user.IsActive,
        user.ID,
    )
    return err
}

// ============ DELETE USER ============
func DeleteUser(id uuid.UUID) error {
    res, err := DB.Exec(`DELETE FROM users WHERE id=$1`, id)
    if err != nil {
        return err
    }
    count, _ := res.RowsAffected()
    if count == 0 {
        return errors.New("user not found")
    }
    return nil
}



