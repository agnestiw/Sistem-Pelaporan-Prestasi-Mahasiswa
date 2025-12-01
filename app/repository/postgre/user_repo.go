package postgre

import (
	"sistem-prestasi/app/model/postgre"
	"sistem-prestasi/database"
)

func FindByUsername(username string) (*postgre.User, error) {
	query := `
		SELECT u.id, u.username, u.password_hash, u.full_name, u.role_id, r.name as role_name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.username = $1 AND u.is_active = true
	`

	var user postgre.User
	err := database.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.RoleName,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetPermissionsByRoleID(roleID string) ([]string, error) {
	query := `
		SELECT p.name 
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`

	rows, err := database.DB.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func UserFindByID(id string) (*postgre.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.full_name, u.role_id, r.name as role_name, u.is_active, u.created_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`

	var user postgre.User
	err := database.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.RoleID,
		&user.RoleName,
		&user.IsActive,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FindAll() ([]postgre.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.full_name, u.role_id, r.name as role_name, u.is_active, u.created_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
		ORDER BY u.created_at DESC
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []postgre.User
	for rows.Next() {
		var u postgre.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.RoleName, &u.IsActive, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func Create(user postgre.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, full_name, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5, true)
	`
	_, err := database.DB.Exec(query, user.Username, user.Email, user.PasswordHash, user.FullName, user.RoleID)
	return err
}

func Update(id string, user postgre.User) error {
	query := `
		UPDATE users 
		SET username = $1, email = $2, full_name = $3, is_active = $4, updated_at = NOW()
		WHERE id = $5
	`
	_, err := database.DB.Exec(query, user.Username, user.Email, user.FullName, user.IsActive, id)
	return err
}

func UpdatePassword(id string, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1 WHERE id = $2`
	_, err := database.DB.Exec(query, passwordHash, id)
	return err
}

func UpdateRole(id string, roleID string) error {
	query := `UPDATE users SET role_id = $1 WHERE id = $2`
	_, err := database.DB.Exec(query, roleID, id)
	return err
}

func Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := database.DB.Exec(query, id)
	return err
}
