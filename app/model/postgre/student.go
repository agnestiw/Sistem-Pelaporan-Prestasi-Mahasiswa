package postgre

import "time"

type Student struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	StudentID    string    `json:"studentId" db:"student_id"` // NIM
	ProgramStudy string    `json:"programStudy" db:"program_study"`
	AcademicYear string    `json:"academicYear" db:"academic_year"`
	AdvisorID    *string   `json:"advisorId" db:"advisor_id"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

type StudentDetail struct {
	ID           string  `json:"id"`
	UserID       string  `json:"userId"`
	StudentID    string  `json:"studentId"`
	FullName     string  `json:"fullName"`
	Email        string  `json:"email"`
	ProgramStudy string  `json:"programStudy"`
	AdvisorName  *string `json:"advisorName"`
}

type AssignAdvisorRequest struct {
	LecturerID string `json:"lecturerId" validate:"required"`
}