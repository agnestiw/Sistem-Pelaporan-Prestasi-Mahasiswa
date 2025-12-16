package service

import (
	repoPg "sistem-prestasi/app/repository/postgre"
	repoMongo "sistem-prestasi/app/repository/mongo"

	"github.com/gofiber/fiber/v2"
)

// @Summary Get achievement statistics
// @Description Get global achievement statistics (admin / authorized)
// @Tags Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/statistics [get]
func GetStatisticsService(c *fiber.Ctx) error {
	totalByStatus, err := repoPg.GetTotalAchievementByStatusRepo()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed get total by status",
			"error":   err.Error(),
		})
	}

	totalByPeriod, err := repoPg.GetTotalAchievementByPeriodRepo()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed get total by period",
			"error":   err.Error(),
		})
	}

	topStudents, err := repoPg.GetTopStudentsRepo()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed get top students",
			"error":   err.Error(),
		})
	}

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

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"total_by_type":            totalByStatus, 
			"total_by_period":          totalByPeriod,
			"top_students":             topStudents,
			"competition_distribution": competitionDistribution,
		},
	})
}

// @Summary Get student achievement report
// @Description Get achievement statistics for a specific student
// @Tags Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/student/{id} [get]
func GetStudentReportService(c *fiber.Ctx) error {
	studentID := c.Params("id")
	if studentID == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "student id is required",
		})
	}

	totalByStatus, err := repoPg.GetStudentTotalByStatusRepo(studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed get total by status",
			"error":   err.Error(),
		})
	}

	totalByPeriod, err := repoPg.GetStudentTotalByPeriodRepo(studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed get total by period",
			"error":   err.Error(),
		})
	}

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

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"total_by_type":            totalByStatus,
			"total_by_period":          totalByPeriod,
			"competition_distribution": competitionDistribution,
		},
	})
}
