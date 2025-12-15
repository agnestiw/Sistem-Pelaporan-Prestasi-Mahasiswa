package postgre

import "time"

type Student struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	StudentID    string    `json:"student_id" db:"student_id"` // NIM
	ProgramStudy string    `json:"programStudy" db:"program_study"`
	AcademicYear string    `json:"academicYear" db:"academic_year"`
	AdvisorID    *string   `json:"advisorId" db:"advisor_id"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

type StudentDetail struct {
	ID           string  `json:"id"`
	UserID       string  `json:"userId"`
	StudentID    string  `json:"student_id"`
	FullName     string  `json:"fullName"`
	Email        string  `json:"email"`
	ProgramStudy string  `json:"programStudy"`
	AdvisorName  *string `json:"advisorName"`
}

type StudentAchievement struct {
	StudentDetail
	AchievementReference
}

type AssignAdvisorRequest struct {
	LecturerID string `json:"lecturerId" validate:"required"`
}