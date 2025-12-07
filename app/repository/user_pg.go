// app/repository/user_pg.go
package repository

import (
    "database/sql" // tambahkan ini
    "projectuas/app/model"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

// HAPUS BARIS INI: var DB *sql.DB
// Biarkan DB di-init dari file lain (achievement_pg.go atau main.go)

func InitDB(db *sql.DB) {
    DB = db // tetap pakai DB dari sini
}

func CreateUser(u *model.User, plainPassword string) error {
    hashed, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), 10)
    u.ID = uuid.New()
    _, err := DB.Exec(`
        INSERT INTO users (id, username, password_hash, full_name, role)
        VALUES ($1, $2, $3, $4, $5)`,
        u.ID, u.Username, string(hashed), u.FullName, u.Role)
    return err
}

func GetUserByUsername(username string) (*model.User, error) {
    u := &model.User{}
    
    err := DB.QueryRow(`
        SELECT id, username, password_hash, full_name, role 
        FROM users 
        WHERE username = $1`, username).
        Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Role)

    if err == sql.ErrNoRows {
        return nil, nil  // INI YANG BENAR!! Kembalikan nil, nil → artinya user tidak ada
    }
    if err != nil {
        return nil, err
    }
    
    return u, nil
}

func GetAllUsers() ([]model.User, error) {
    rows, err := DB.Query(`SELECT id, username, full_name, role FROM users`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var users []model.User
    for rows.Next() {
        var u model.User
        if err := rows.Scan(&u.ID, &u.Username, &u.FullName, &u.Role); err != nil {
            return nil, err
        }
        users = append(users, u)
    }
    return users, nil
}

func GetUserByID(idStr string) (*model.User, error) {
    id, err := uuid.Parse(idStr)
    if err != nil {
        return nil, err
    }
    u := &model.User{}
    err = DB.QueryRow(`
        SELECT id, username, full_name, role 
        FROM users WHERE id = $1`, id).
        Scan(&u.ID, &u.Username, &u.FullName, &u.Role)
    return u, err
}

func UpdateUserRole(id uuid.UUID, role string) error {
    _, err := DB.Exec(`UPDATE users SET role = $1 WHERE id = $2`, role, id)
    return err
}

func DeleteUser(id uuid.UUID) error {
    _, err := DB.Exec(`DELETE FROM users WHERE id = $1`, id)
    return err
}

