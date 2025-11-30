package postgre

import "time"

type Lecturer struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	LecturerID string    `json:"lecturerId" db:"lecturer_id"` // NIP
	Department string    `json:"department" db:"department"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}

type LecturerDetail struct {
	ID         string `json:"id"`
	LecturerID string `json:"lecturerId"` // NIP
	FullName   string `json:"fullName"`
	Email      string `json:"email"`
	Department string `json:"department"`
}