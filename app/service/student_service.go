package service

import (
	"context"
	repoMongo "sistem-prestasi/app/repository/mongo"
	repoPostgre "sistem-prestasi/app/repository/postgre"

	"github.com/gofiber/fiber/v2"
)


// @Summary Get all students
// @Description Mengambil semua data mahasiswa (hanya admin & dosen)
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/students [get]
func GetAllStudentService(c *fiber.Ctx) error {

	nama_role := c.Locals("role_name")

	if nama_role == "Mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"message": "anda bukan seorang admin maupun dosen",
		})
	}

	students, err := repoPostgre.GetAllStudentRepo()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil data mahasiswa",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   students,
	})
}


// @Summary Get student by ID
// @Description Mengambil data mahasiswa berdasarkan ID
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/students/{id} [get]
func GetStudentByID(c *fiber.Ctx) error {

	id := c.Params("id")
	if id == "" {
		return c.Status(404).JSON(fiber.Map{
			"message": "student id tidak valid",
		})
	}

	hasil, err := repoPostgre.GetStudentByIDRepo(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "tidak dapat mengambil data student",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   hasil,
	})

}


// @Summary Get student achievement detail
// @Description Mengambil detail prestasi mahasiswa beserta data achievement dari MongoDB
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/students/{id}/achievements [get]
func GetStudentAchievementDetailService(c *fiber.Ctx) error {

	id := c.Params("id")
	if id == "" {
		return c.Status(404).JSON(fiber.Map{
			"message": "id student tidak valid",
		})
	}

	result, err := repoPostgre.GetStudentAchievementDetailRepo(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil data reference",
			"error":   err.Error(),
		})
	}

	mongoData, err := repoMongo.FindAchievementByID(context.Background(), result.MongoAchievementID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil data achievement di MongoDB",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"student": result.StudentDetail,
		"reference": result.AchievementReference,
		"achievement": mongoData,
	})

}


// @Summary Set student advisor
// @Description Menentukan dosen pembimbing (advisor) untuk mahasiswa
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Param body body map[string]string true "Advisor ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/students/{id}/advisor [put]
func SetStudentAdvisorService(c *fiber.Ctx) error {

	studentID := c.Params("id")
	if studentID == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "student_id tidak valid",
		})
	}

	// ambil advisor_id dari body
	var body map[string]string
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "body tidak valid",
			"error":   err.Error(),
		})
	}

	advisorID := body["advisor_id"]
	if advisorID == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "advisor_id wajib diisi",
		})
	}

	result, err := repoPostgre.SetStudentAdvisorRepo(studentID, advisorID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "tidak dapat set advisor ke student",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

