package postgre

import "time"

type User struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FullName     string    `json:"full_name" db:"full_name"`
	RoleID       string    `json:"role_id" db:"role_id"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	
	// Helper field untuk join query (tidak ada di tabel users, tapi berguna)
	RoleName     string    `json:"role_name,omitempty" db:"role_name"`
}