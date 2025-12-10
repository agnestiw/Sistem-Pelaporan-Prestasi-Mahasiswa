package postgre

import (
	"database/sql"
	modelPostgre "sistem-prestasi/app/model/postgre"
	"sistem-prestasi/database"
)

func FindStudentByUserID(userID string) (*modelPostgre.StudentDetail, error) {
    query := `
        SELECT s.id, s.user_id, s.student_id, u.full_name, u.email, s.program_study, u_lec.full_name as advisor_name
        FROM students s
        JOIN users u ON s.user_id = u.id
        LEFT JOIN lecturers l ON s.advisor_id = l.id
        LEFT JOIN users u_lec ON l.user_id = u_lec.id
        WHERE s.advisor_id = $1
    `
    var s modelPostgre.StudentDetail
    var advisorName sql.NullString

    err := database.DB.QueryRow(query, userID).Scan(&s.ID, &s.UserID, &s.StudentID, &s.FullName, &s.Email, &s.ProgramStudy, &advisorName)
    if err != nil {
        return nil, err
    }

    if advisorName.Valid {
        name := advisorName.String
        s.AdvisorName = &name
    }
    return &s, nil
}


func StudentFindAll() ([]modelPostgre.StudentDetail, error) {
	query := `
		SELECT s.id, s.user_id, s.student_id, u.full_name, u.email, s.program_study, u_lec.full_name as advisor_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		LEFT JOIN lecturers l ON s.advisor_id = l.id
		LEFT JOIN users u_lec ON l.user_id = u_lec.id
		ORDER BY s.student_id ASC
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []modelPostgre.StudentDetail
	for rows.Next() {
		var s modelPostgre.StudentDetail
		var advisorName sql.NullString

		err := rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.FullName, &s.Email, &s.ProgramStudy, &advisorName)
		if err != nil {
			return nil, err
		}

		if advisorName.Valid {
			name := advisorName.String
			s.AdvisorName = &name
		}
		students = append(students, s)
	}
	return students, nil
}

func StudentFindByID(id string) (*modelPostgre.StudentDetail, error) {
	query := `
		SELECT s.id, s.user_id, s.student_id, u.full_name, u.email, s.program_study, u_lec.full_name as advisor_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		LEFT JOIN lecturers l ON s.advisor_id = l.id
		LEFT JOIN users u_lec ON l.user_id = u_lec.id
		WHERE s.id = $1
	`
	var s modelPostgre.StudentDetail
	var advisorName sql.NullString

	err := database.DB.QueryRow(query, id).Scan(&s.ID, &s.UserID, &s.StudentID, &s.FullName, &s.Email, &s.ProgramStudy, &advisorName)
	if err != nil {
		return nil, err
	}

	if advisorName.Valid {
		name := advisorName.String
		s.AdvisorName = &name
	}
	return &s, nil
}

func GetStudentByUserID(userID string) (string, error) {
    var studentID string
    query := `SELECT id FROM students WHERE user_id = $1`
    
    err := database.DB.QueryRow(query, userID).Scan(&studentID)
    if err != nil {
        return "", err
    }
    
    return studentID, nil
}

func UpdateAdvisor(studentID, lecturerID string) error {
	query := `UPDATE students SET advisor_id = $1 WHERE id = $2`
	_, err := database.DB.Exec(query, lecturerID, studentID)
	return err
}
