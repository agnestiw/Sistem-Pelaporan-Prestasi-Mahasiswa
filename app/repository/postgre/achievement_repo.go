package postgre

import (
	modelPostgres "sistem-prestasi/app/model/postgre"
	// modelMongo "sistem-prestasi/app/model/postgre"
	"sistem-prestasi/database"
	"time"
)

func GetStudentIdFromAchievementReferences(achievementReferenceID string) (string, error) {
	var studentID string

	query := `
		SELECT student_id
		FROM achievement_references
		WHERE id = $1
	`

	err := database.DB.QueryRow(query, achievementReferenceID).Scan(&studentID)
	if err != nil {
		return "", err
	}

	return studentID, nil
}


func GetAllAchievementsRepo() ([]modelPostgres.AchievementReference, error) {
	query := `
        SELECT 
            ar.id, 
            ar.student_id, 
            ar.mongo_achievement_id, 
            ar.status, 
            ar.submitted_at, 
            ar.verified_at, 
            ar.verified_by, 
            ar.rejection_note, 
            ar.created_at
        FROM achievement_references ar
        ORDER BY ar.created_at DESC
    `

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []modelPostgres.AchievementReference

	for rows.Next() {
		var data modelPostgres.AchievementReference

		err := rows.Scan(
			&data.ID,
			&data.StudentID,
			&data.MongoAchievementID,
			&data.Status,
			&data.SubmittedAt,
			&data.VerifiedAt,
			&data.VerifiedBy,
			&data.RejectionNote,
			&data.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, data)
	}

	return results, nil
}

func CreateAchievementRef(ref modelPostgres.AchievementReference) error {
	query := `
		INSERT INTO achievement_references (
			id, student_id, mongo_achievement_id, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := database.DB.Exec(query,
		ref.ID, ref.StudentID, ref.MongoAchievementID, ref.Status, ref.CreatedAt, ref.UpdatedAt)
	return err
}


func GetAllAchievementByStudentID(studentID string) ([]modelPostgres.AchievementReference, error) {
    query := `
        SELECT 
            ar.id, 
            ar.student_id, 
            ar.mongo_achievement_id, 
            ar.status, 
            ar.submitted_at, 
            ar.verified_at, 
            ar.verified_by, 
            ar.rejection_note, 
            ar.created_at
        FROM achievement_references ar
        WHERE ar.student_id = $1
        ORDER BY ar.created_at DESC
    `

    rows, err := database.DB.Query(query, studentID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []modelPostgres.AchievementReference

    for rows.Next() {
        var data modelPostgres.AchievementReference

        err := rows.Scan(
            &data.ID,
            &data.StudentID,
            &data.MongoAchievementID,
            &data.Status,
            &data.SubmittedAt,
            &data.VerifiedAt,
            &data.VerifiedBy,
            &data.RejectionNote,
            &data.CreatedAt,
        )
        if err != nil {
            return nil, err
        }

        results = append(results, data)
    }

    return results, nil
}



func UpdateAchievementRefUpdatedAt(id string) error {
	query := `
		UPDATE achievement_references
		SET updated_at = $1
		WHERE id = $2
	`

	_, err := database.DB.Exec(query, time.Now(), id)
	return err
}

func GetAchievementRefByID(id string) (*modelPostgres.AchievementRefWithStudent, error) {

	query := `
        SELECT 
            ar.id, 
            ar.student_id, 
            ar.mongo_achievement_id, 
            ar.status, 
            ar.submitted_at, 
            ar.verified_at,      -- Tambahkan ini agar lengkap
            ar.verified_by,      -- Tambahkan ini agar lengkap
            ar.rejection_note,   -- Tambahkan ini agar lengkap
            ar.created_at,
            u.full_name as student_name
        FROM achievement_references ar
        JOIN students s ON ar.student_id = s.id
        JOIN users u ON s.user_id = u.id
        WHERE ar.id = $1
    `

	var data modelPostgres.AchievementRefWithStudent

	err := database.DB.QueryRow(query, id).Scan(
		&data.ID,
		&data.StudentID,
		&data.MongoAchievementID,
		&data.Status,
		&data.SubmittedAt,
		&data.VerifiedAt,
		&data.VerifiedBy,
		&data.RejectionNote,
		&data.CreatedAt,
		&data.StudentName,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func SubmitAchievementRepo(achievement_references_id string) (bool, error) {

	query, err := database.DB.Exec(`
		UPDATE achievement_references
		SET status = 'submitted',
			submitted_at = NOW()
		WHERE id = $1
    `, achievement_references_id)

	if err != nil {
		return false, err
	}

	rowsEffected, err := query.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsEffected == 0 {
		return false, nil
	}

	return true, err

}
