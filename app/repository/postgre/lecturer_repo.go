package postgre

import (
	"database/sql"
	modelPostgre "sistem-prestasi/app/model/postgre"
	"sistem-prestasi/database"
)

func FindLecturerByUserID(userID string) (*modelPostgre.Lecturer, error) {
    var l modelPostgre.Lecturer
    query := `SELECT id, user_id, lecturer_id, department, created_at FROM lecturers WHERE user_id = $1`
    err := database.DB.QueryRow(query, userID).Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &l, nil
}

func FindAllLecturers() ([]modelPostgre.LecturerDetail, error) {
	query := `
        SELECT l.id, l.lecturer_id, u.full_name, u.email, l.department
        FROM lecturers l
        JOIN users u ON l.user_id = u.id
        ORDER BY u.full_name ASC
    `
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lecturers []modelPostgre.LecturerDetail
	for rows.Next() {
		var l modelPostgre.LecturerDetail
		if err := rows.Scan(&l.ID, &l.LecturerID, &l.FullName, &l.Email, &l.Department); err != nil {
			return nil, err
		}
		lecturers = append(lecturers, l)
	}
	return lecturers, nil
}

func FindLecturerByID(id string) (*modelPostgre.Lecturer, error) {
	var l modelPostgre.Lecturer

	query := `SELECT id, user_id, lecturer_id, department, created_at FROM lecturers WHERE id = $1`
	err := database.DB.QueryRow(query, id).Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func FindLecturerAdvisees(lecturerID string) ([]modelPostgre.StudentDetail, error) {
	query := `
        SELECT s.id, s.user_id, s.student_id, u.full_name, u.email, s.program_study, u_lec.full_name as advisor_name
        FROM students s
        JOIN users u ON s.user_id = u.id
        JOIN lecturers l ON s.advisor_id = l.id
        JOIN users u_lec ON l.user_id = u_lec.id
        WHERE s.advisor_id = $1
    `
	rows, err := database.DB.Query(query, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []modelPostgre.StudentDetail
	for rows.Next() {
		var s modelPostgre.StudentDetail
		var advisorName sql.NullString

		if err := rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.FullName, &s.Email, &s.ProgramStudy, &advisorName); err != nil {
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


func GetLecturersRepo() ([]modelPostgre.LecturerDetail, error) {
	var lecturers []modelPostgre.LecturerDetail
	
	rows, err := database.DB.Query(`
		SELECT 
			l.id,
			l.lecturer_id,
			u.full_name,
			u.email,
			l.department
		FROM lecturers l
		JOIN users u ON u.id = l.user_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var lecturer modelPostgre.LecturerDetail
		err := rows.Scan(
			&lecturer.ID,
			&lecturer.LecturerID,
			&lecturer.FullName,
			&lecturer.Email,
			&lecturer.Department,
		)
		if err != nil {
			return nil, err
		}
		lecturers = append(lecturers, lecturer)
	}

	return lecturers, nil
}
