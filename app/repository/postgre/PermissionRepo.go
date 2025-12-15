package postgre

import "sistem-prestasi/database"

func UserHasPermissionRepo(userID string, permission string) (bool, error) {

	var exists bool

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users u
			JOIN role_permissions rp ON u.role_id = rp.role_id
			JOIN permissions p ON rp.permission_id = p.id
			WHERE u.id = $1
			AND p.name = $2
		)
	`

	err := database.DB.QueryRow(query, userID, permission).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
