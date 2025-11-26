package postgre

import "time"

type Student struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	StudentID    string    `json:"student_id" db:"student_id"` // NIM
	ProgramStudy string    `json:"program_study" db:"program_study"`
	AcademicYear string    `json:"academic_year" db:"academic_year"`
	AdvisorID    *string   `json:"advisor_id" db:"advisor_id"` // Pointer = Nullable
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	
	// Relation Helper
	AdvisorName  string    `json:"advisor_name,omitempty" db:"advisor_name"`
}