package postgre

import (
	"database/sql"
	"errors"
	"sistem-prestasi/app/model/postgre"
	"sistem-prestasi/database"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(raw string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
	return err == nil
}

func GetProfile(userId string) (*postgre.User, error) {

	var User postgre.User

	err := database.DB.QueryRow(`
		SELECT 
			u.id, u.username, u.email, u.full_name, u.password_hash, u.role_id, r.name  
		FROM 
			users as u
		JOIN 
			roles as r on u.role_id = r.id
		WHERE 
			u.id = $1 
	`, userId).Scan(
		&User.ID, &User.Username, &User.Email, &User.FullName, &User.PasswordHash, &User.RoleID, &User.RoleName,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("tidak ditemukan")
	}

	return &User, err

}

func Authenticate(username, password string) (*postgre.User, error) {
	var user postgre.User
	var studentID sql.NullString
	var advisorID sql.NullString

	err := database.DB.QueryRow(`
		SELECT 
			u.id,
			u.username,
			u.email,
			u.full_name,
			u.password_hash,
			u.role_id,
			r.name,
			s.id AS student_id,
			l.id AS lecturer_id
		FROM users u
		JOIN roles r ON u.role_id = r.id
		LEFT JOIN students s ON s.user_id = u.id
		LEFT JOIN lecturers l ON l.user_id = u.id
		WHERE u.username = $1
	`, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.PasswordHash,
		&user.RoleID,
		&user.RoleName,
		&studentID,
		&advisorID,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("username tidak ditemukan")
	}
	if err != nil {
		// ⬅️ INI YANG MEMBONGKAR MASALAH ASLI
		return nil, err
	}

	if !CheckPassword(password, user.PasswordHash) {
		return nil, errors.New("password salah")
	}

	if studentID.Valid {
		user.StudentID = &studentID.String
	}
	if advisorID.Valid {
		user.AdvisorID = &advisorID.String
	}

	return &user, nil
}

