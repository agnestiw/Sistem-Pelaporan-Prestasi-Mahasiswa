package postgre

import (
	modelPostgres "sistem-prestasi/app/model/postgre"
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

func GetAdvisorIDByAchievementRef(refID string) (string, error) {
	var advisorID string

	err := database.DB.QueryRow(`
		SELECT s.advisor_id
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		WHERE ar.id = $1
	`, refID).Scan(&advisorID)

	return advisorID, err
}

func GetAllAchievementsRepo(limit, offset int) ([]modelPostgres.AchievementReference, int, error) {
    var total int
    // Hitung total data untuk pagination metadata
    countQuery := `SELECT COUNT(*) FROM achievement_references WHERE status != 'deleted'`
    err := database.DB.QueryRow(countQuery).Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    query := `
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
               verified_at, verified_by, rejection_note, created_at
        FROM achievement_references
        WHERE status != 'deleted'
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
    rows, err := database.DB.Query(query, limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var results []modelPostgres.AchievementReference
    for rows.Next() {
        var data modelPostgres.AchievementReference
        err := rows.Scan(&data.ID, &data.StudentID, &data.MongoAchievementID, &data.Status, 
                         &data.SubmittedAt, &data.VerifiedAt, &data.VerifiedBy, 
                         &data.RejectionNote, &data.CreatedAt)
        if err != nil {
            return nil, 0, err
        }
        results = append(results, data)
    }

    return results, total, nil
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

func GetAchievementDetailByAchievementIDRepo(achievementID string,) (*modelPostgres.AchievementRefWithStudent, error) {
	var result modelPostgres.AchievementRefWithStudent
	err := database.DB.QueryRow(`
		SELECT
			ar.id,
			ar.student_id,
			ar.mongo_achievement_id,
			ar.status,
			ar.submitted_at,
			ar.verified_at,
			ar.verified_by,
			ar.rejection_note,
			ar.created_at,
			ar.updated_at,
			u.full_name AS student_name
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		JOIN users u ON u.id = s.user_id
		WHERE ar.id = $1
	`, achievementID).Scan(
		&result.ID,
		&result.StudentID,
		&result.MongoAchievementID,
		&result.Status,
		&result.SubmittedAt,
		&result.VerifiedAt,
		&result.VerifiedBy,
		&result.RejectionNote,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.StudentName,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}


func GetAllAchievementByStudentID(studentID string, limit, offset int) ([]modelPostgres.AchievementReference, int, error) {
    var total int
    countQuery := `SELECT COUNT(*) FROM achievement_references WHERE student_id = $1 AND status != 'deleted'`
    err := database.DB.QueryRow(countQuery, studentID).Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    query := `
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
               verified_at, verified_by, rejection_note, created_at
        FROM achievement_references
        WHERE student_id = $1 AND status != 'deleted'
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `
    rows, err := database.DB.Query(query, studentID, limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var results []modelPostgres.AchievementReference
    for rows.Next() {
        var data modelPostgres.AchievementReference
        err := rows.Scan(&data.ID, &data.StudentID, &data.MongoAchievementID, &data.Status, 
                         &data.SubmittedAt, &data.VerifiedAt, &data.VerifiedBy, 
                         &data.RejectionNote, &data.CreatedAt)
        if err != nil {
            return nil, 0, err
        }
        results = append(results, data)
    }

    return results, total, nil
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

func SoftDeleteAchievementRef(id string) error {
	query := `
		UPDATE achievement_references
		SET status = 'deleted', updated_at = NOW()
		WHERE id = $1
	`

	_, err := database.DB.Exec(query, id)
	return err
}


func GetAchievementRefByID(id string) (*modelPostgres.AchievementRefWithStudent, error) {
	var data modelPostgres.AchievementRefWithStudent
	
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
            ar.created_at,
            u.full_name as student_name
        FROM achievement_references ar
        JOIN students s ON ar.student_id = s.id
        JOIN users u ON s.user_id = u.id
        WHERE ar.id = $1
    `

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

func VerifyAchievementRepo(achievement_references_id string) (bool, error) {

	query, err := database.DB.Exec(`
		UPDATE achievement_references
		SET status = 'verified',
			verified_at = NOW()
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

func RejectAchievementRepo(achievement_references_id string, rejection_note string, verified_by string) (bool, error) {

	query, err := database.DB.Exec(`
		UPDATE achievement_references
		SET status = 'rejected',
			rejection_note = $2,
			verified_by = $3,
			verified_at = NOW()
		WHERE id = $1
	`, achievement_references_id, rejection_note, verified_by)

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
