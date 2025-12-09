package repository

import (
    "context"
    "database/sql"
    "os"
    "time"

    "projectuas/app/model"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

var DB *sql.DB
var MongoClient *mongo.Client

// InitDB dipanggil dari route/setup.go atau main.go
func Init(pg *sql.DB, mongo *mongo.Client) {
    DB = pg
    MongoClient = mongo
}

// ------------------------
// USER
// ------------------------
func CreateUser(u *model.User, plainPassword string) error {
    hashed, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), 10)
    u.ID = uuid.New()
    _, err := DB.Exec(`
        INSERT INTO users (id, username, password_hash, full_name, role_id)
        VALUES ($1, $2, $3, $4, $5)`,
        u.ID, u.Username, string(hashed), u.FullName, u.RoleID)
    return err
}

func GetUserByUsername(username string) (*model.User, error) {
    u := &model.User{}
    err := DB.QueryRow(`
        SELECT id, username, password_hash, full_name, role_id
        FROM users WHERE username = $1`, username).
        Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.RoleID)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return u, err
}

func GetAllUsers() ([]model.User, error) {
    rows, err := DB.Query(`SELECT id, username, full_name, role_id FROM users`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []model.User
    for rows.Next() {
        var u model.User
        if err := rows.Scan(&u.ID, &u.Username, &u.FullName, &u.RoleID); err != nil {
            return nil, err
        }
        users = append(users, u)
    }
    return users, nil
}

func GetUserByID(id uuid.UUID) (*model.User, error) {
    u := &model.User{}
    err := DB.QueryRow(`SELECT id, username, full_name, role_id FROM users WHERE id=$1`, id).
        Scan(&u.ID, &u.Username, &u.FullName, &u.RoleID)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return u, err
}

func UpdateUserRole(id uuid.UUID, roleID uuid.UUID) error {
    _, err := DB.Exec(`UPDATE users SET role_id=$1 WHERE id=$2`, roleID, id)
    return err
}

func DeleteUser(id uuid.UUID) error {
    _, err := DB.Exec(`DELETE FROM users WHERE id=$1`, id)
    return err
}

// ------------------------
// ROLE
// ------------------------
func GetRoleByName(name string) (*model.Role, error) {
    role := &model.Role{}
    err := DB.QueryRow(`SELECT id, name FROM roles WHERE name=$1`, name).Scan(&role.ID, &role.Name)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return role, err
}

func GetRoleByID(id uuid.UUID) (*model.Role, error) {
    role := &model.Role{}
    err := DB.QueryRow(`SELECT id, name FROM roles WHERE id=$1`, id).Scan(&role.ID, &role.Name)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return role, err
}

// ------------------------
// ACHIEVEMENT (PostgreSQL)
// ------------------------
func GetRef(id uuid.UUID) (*model.AchievementRef, error) {
    ref := &model.AchievementRef{}
    err := DB.QueryRow(`
        SELECT id, student_id, mongo_achievement_id, status, submitted_at,
               verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references WHERE id = $1`, id).
        Scan(&ref.ID, &ref.StudentID, &ref.MongoID, &ref.Status,
            &ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy,
            &ref.RejectionNote, &ref.CreatedAt, &ref.UpdatedAt)
    return ref, err
}

func SubmitAchievementPG(id uuid.UUID) error {
    _, err := DB.Exec(`UPDATE achievement_references SET status='submitted', submitted_at=NOW(), updated_at=NOW() WHERE id=$1`, id)
    return err
}

func UpdateToDeleted(id uuid.UUID) error {
    _, err := DB.Exec(`UPDATE achievement_references SET status='deleted', updated_at=NOW() WHERE id=$1`, id)
    return err
}

// ------------------------
// ACHIEVEMENT (MongoDB Soft Delete)
// ------------------------
func SoftDeleteMongo(mongoID string) error {
    coll := MongoClient.Database(os.Getenv("DB_NAME")).Collection("achievements")
    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    _, err := coll.UpdateOne(ctx, bson.M{"_id": mongoID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
    return err
}
