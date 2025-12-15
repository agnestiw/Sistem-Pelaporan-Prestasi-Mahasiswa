package service

import (
	repoPg "sistem-prestasi/app/repository/postgre"
	repoMongo "sistem-prestasi/app/repository/mongo"

	"github.com/gofiber/fiber/v2"
)




func GetStatisticsService(c *fiber.Ctx) error {

	// 1️⃣ Total prestasi per STATUS (PostgreSQL)
	totalByStatus, err := repoPg.GetTotalAchievementByStatusRepo()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed get total by status",
			"error":   err.Error(),
		})
	}

	// 2️⃣ Total prestasi per periode (PostgreSQL)
	totalByPeriod, err := repoPg.GetTotalAchievementByPeriodRepo()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed get total by period",
			"error":   err.Error(),
		})
	}

	// 3️⃣ Top mahasiswa berprestasi (PostgreSQL)
	topStudents, err := repoPg.GetTopStudentsRepo()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed get top students",
			"error":   err.Error(),
		})
	}

	// 4️⃣ Distribusi tingkat kompetisi (MongoDB)
	mongoIDs, err := repoPg.GetVerifiedCompetitionMongoIDsRepo()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed get mongo ids",
			"error":   err.Error(),
		})
	}

	competitionDistribution, err := repoMongo.GetCompetitionLevelDistributionMongo(mongoIDs)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed get competition distribution",
			"error":   err.Error(),
		})
	}

	// ✅ RESPONSE FINAL (LENGKAP)
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"total_by_type":            totalByStatus,   // status / atau type jika sudah ada
			"total_by_period":          totalByPeriod,
			"top_students":             topStudents,
			"competition_distribution": competitionDistribution,
		},
	})
}


func GetStudentReportService(c *fiber.Ctx) error {

	studentID := c.Params("id")
	if studentID == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "student id is required",
		})
	}

	// 1️⃣ Total prestasi per status
	totalByStatus, err := repoPg.GetStudentTotalByStatusRepo(studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed get total by status",
			"error":   err.Error(),
		})
	}

	// 2️⃣ Total prestasi per periode
	totalByPeriod, err := repoPg.GetStudentTotalByPeriodRepo(studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed get total by period",
			"error":   err.Error(),
		})
	}

	// 3️⃣ Distribusi tingkat kompetisi (MongoDB)
	mongoIDs, err := repoPg.GetStudentVerifiedMongoIDsRepo(studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed get mongo ids",
			"error":   err.Error(),
		})
	}

	competitionDistribution := []map[string]interface{}{}
	if len(mongoIDs) > 0 {
		competitionDistribution, err = repoMongo.GetCompetitionLevelDistributionMongo(mongoIDs)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "failed get competition distribution",
				"error":   err.Error(),
			})
		}
	}

	// ✅ RESPONSE FINAL
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"total_by_type":            totalByStatus,
			"total_by_period":          totalByPeriod,
			"competition_distribution": competitionDistribution,
		},
	})
}
