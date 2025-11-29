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
	RoleName     string    `json:"role_name,omitempty" db:"role_name"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type LoginResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refreshToken"`
	User         UserDetail  `json:"user"`
}

type UserDetail struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	FullName    string   `json:"fullName"`
	RoleID       string    `json:"role_id" db:"role_id"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=4"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"fullName" validate:"required"`
	RoleID   string `json:"roleId" validate:"required"` // Admin harus pilih role saat create
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" validate:"email"`
	FullName string `json:"fullName"`
	Password string `json:"password"` // Opsional, jika kosong berarti tidak ganti password
	IsActive *bool  `json:"isActive"` // Pointer agar bisa mendeteksi false
}

type AssignRoleRequest struct {
	RoleID string `json:"roleId" validate:"required"`
}