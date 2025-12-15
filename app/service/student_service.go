package service

import (
	"context"
	repoMongo "sistem-prestasi/app/repository/mongo"
	repoPostgre "sistem-prestasi/app/repository/postgre"

	"github.com/gofiber/fiber/v2"
)

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

