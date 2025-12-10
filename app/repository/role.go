package repository

import (
    "projectuas/app/model"
    "github.com/google/uuid"
)

func GetRoleByID(id uuid.UUID) (*model.Role, error) {
    var role model.Role
    query := `SELECT id, name, description FROM roles WHERE id = $1`
    err := DB.QueryRow(query, id).Scan(
        &role.ID,
        &role.Name,
        &role.Description,
    )
    if err != nil {
        return nil, err
    }
    return &role, nil
}

func GetRoleByName(name string) (*model.Role, error) {
    var role model.Role
    query := `SELECT id, name, description FROM roles WHERE name = $1`
    err := DB.QueryRow(query, name).Scan(
        &role.ID,
        &role.Name,
        &role.Description,
    )
    if err != nil {
        return nil, err
    }
    return &role, nil
}

func GetPermissionsByRoleID(roleID uuid.UUID) ([]string, error) {
    query := `
        SELECT p.name
        FROM permissions p
        JOIN role_permissions rp ON p.id = rp.permission_id
        WHERE rp.role_id = $1
    `
    rows, err := DB.Query(query, roleID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var permissions []string
    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            return nil, err
        }
        permissions = append(permissions, name)
    }
    return permissions, nil
}
