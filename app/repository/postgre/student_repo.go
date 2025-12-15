package postgre

import (
	model "sistem-prestasi/app/model/postgre"
	"sistem-prestasi/database"
)

func GetAllStudentRepo() ([]model.Student, error) {
	rows, err := database.DB.Query(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student

	for rows.Next() {
		var s model.Student
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.StudentID,
			&s.ProgramStudy,
			&s.AcademicYear,
			&s.AdvisorID,
			&s.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		students = append(students, s)
	}

	return students, nil
}

func GetStudentByIDRepo(ID string) (model.StudentDetail, error) {

	var student model.StudentDetail

	err := database.DB.QueryRow(`
		SELECT 
			s.id,
			s.user_id,
			s.student_id,
			u.full_name AS student_name,
			u.email,
			s.program_study,

			ua.full_name AS advisor_name
		FROM students s
		JOIN users u ON u.id = s.user_id
		LEFT JOIN lecturers l ON l.id = s.advisor_id
		LEFT JOIN users ua ON ua.id = l.user_id
		WHERE s.id = $1;
	`, ID).Scan(
		&student.ID,
		&student.UserID,
		&student.StudentID,
		&student.FullName,
		&student.Email,
		&student.ProgramStudy,
		&student.AdvisorName,
	)

	if err != nil {
		return model.StudentDetail{}, err
	}

	return student, err

}

func GetStudentAchievementDetailRepo(refID string) (*model.StudentAchievement, error) {
	var result model.StudentAchievement

	err := database.DB.QueryRow(`
		SELECT
			s.id,
			s.user_id,
			s.student_id,
			u.full_name,
			u.email,
			s.program_study,
			ua.full_name AS advisor_name,

			ar.id,
			ar.student_id,
			ar.mongo_achievement_id,
			ar.status,
			ar.submitted_at,
			ar.verified_at,
			ar.verified_by,
			ar.rejection_note,
			ar.created_at,
			ar.updated_at

		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		JOIN users u ON u.id = s.user_id
		LEFT JOIN lecturers l ON l.id = s.advisor_id
		LEFT JOIN users ua ON ua.id = l.user_id
		WHERE s.id = $1
	`, refID).Scan(
		&result.StudentDetail.ID,
		&result.StudentDetail.UserID,
		&result.StudentDetail.StudentID,
		&result.StudentDetail.FullName,
		&result.StudentDetail.Email,
		&result.StudentDetail.ProgramStudy,
		&result.StudentDetail.AdvisorName,

		&result.AchievementReference.ID,
		&result.AchievementReference.StudentID,
		&result.AchievementReference.MongoAchievementID,
		&result.AchievementReference.Status,
		&result.AchievementReference.SubmittedAt,
		&result.AchievementReference.VerifiedAt,
		&result.AchievementReference.VerifiedBy,
		&result.AchievementReference.RejectionNote,
		&result.AchievementReference.CreatedAt,
		&result.AchievementReference.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func SetStudentAdvisorRepo(student_id string, advisor_id string) (bool, error) {

	query, err := database.DB.Exec(`
		UPDATE students
			SET advisor_id = $2
		WHERE id = $1
	`, student_id, advisor_id)

	if err != nil {
		return false, err
	}

	result, err := query.RowsAffected()
	if err != nil {
		return false, err
	}

	if result == 0 {
		return false, nil
	}

	return true, err


}
