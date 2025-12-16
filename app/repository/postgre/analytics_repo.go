package postgre

import (
	"sistem-prestasi/database"
)

func GetTotalAchievementByStatusRepo() ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`
		SELECT status, COUNT(*) AS total
		FROM achievement_references
		GROUP BY status
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var status string
		var total int
		rows.Scan(&status, &total)

		result = append(result, map[string]interface{}{
			"status": status,
			"total":  total,
		})
	}

	return result, nil
}


func GetTotalAchievementByPeriodRepo() ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`
		SELECT DATE_TRUNC('month', verified_at) AS period, COUNT(*) AS total
		FROM achievement_references
		WHERE status = 'verified'
		GROUP BY period
		ORDER BY period
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var period string
		var total int
		rows.Scan(&period, &total)

		result = append(result, map[string]interface{}{
			"period": period,
			"total":  total,
		})
	}

	return result, nil
}


func GetTopStudentsRepo() ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`
		SELECT s.id, u.full_name, COUNT(*) AS total
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		JOIN users u ON s.user_id = u.id
		WHERE ar.status = 'verified'
		GROUP BY s.id, u.full_name
		ORDER BY total DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, name string
		var total int
		rows.Scan(&id, &name, &total)

		result = append(result, map[string]interface{}{
			"student_id": id,
			"name":       name,
			"total":      total,
		})
	}

	return result, nil
}

func GetVerifiedCompetitionMongoIDsRepo() ([]string, error) {
	rows, err := database.DB.Query(`
		SELECT mongo_achievement_id
		FROM achievement_references
		WHERE status = 'verified'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids, nil
}



func GetStudentStatisticsRepo(studentID string) ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`
		SELECT a.achievement_type,
		       a.details->>'level' AS level,
		       ar.verified_at
		FROM achievement_references ar
		JOIN achievements a ON ar.mongo_achievement_id = a.mongo_id
		WHERE ar.student_id = $1
		AND ar.status = 'verified'
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var t, level, verifiedAt string
		rows.Scan(&t, &level, &verifiedAt)

		result = append(result, map[string]interface{}{
			"type":        t,
			"level":       level,
			"verified_at": verifiedAt,
		})
	}

	return result, nil
}



func GetStudentTotalByStatusRepo(studentID string) ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`
		SELECT status, COUNT(*) AS total
		FROM achievement_references
		WHERE student_id = $1
		GROUP BY status
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var status string
		var total int
		rows.Scan(&status, &total)

		result = append(result, map[string]interface{}{
			"status": status,
			"total":  total,
		})
	}

	return result, nil
}

func GetStudentTotalByPeriodRepo(studentID string) ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`
		SELECT DATE_TRUNC('month', verified_at) AS period, COUNT(*) AS total
		FROM achievement_references
		WHERE student_id = $1
		AND status = 'verified'
		GROUP BY period
		ORDER BY period
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var period string
		var total int
		rows.Scan(&period, &total)

		result = append(result, map[string]interface{}{
			"period": period,
			"total":  total,
		})
	}

	return result, nil
}

func GetStudentVerifiedMongoIDsRepo(studentID string) ([]string, error) {
	rows, err := database.DB.Query(`
		SELECT mongo_achievement_id
		FROM achievement_references
		WHERE student_id = $1
		AND status = 'verified'
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids, nil
}