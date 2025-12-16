package service

import (
	repoPostgre "sistem-prestasi/app/repository/postgre"
	"sistem-prestasi/helper"

	"github.com/gofiber/fiber/v2"
)


func GetMyAdvisor(c *fiber.Ctx) error {
	loggedInUserID := c.Locals("user_id").(string)

	student, err := repoPostgre.GetStudentByIDRepo(loggedInUserID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Data mahasiswa tidak ditemukan. Apakah anda login sebagai Mahasiswa?",
		})
	}

	if student.AdvisorName == nil {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Anda belum memiliki Dosen Wali",
			"data":    nil,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"advisorName": *student.AdvisorName,
		},
	})
}


func GetAllLecturers(c *fiber.Ctx) error {
	lecturers, err := repoPostgre.FindAllLecturers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data dosen"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "data": lecturers})
}

func GetLecturerAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")

	lecturer, err := repoPostgre.FindLecturerByID(lecturerID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Dosen tidak ditemukan"})
	}

	loggedInUserID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
	}

	if !helper.IsAdmin(c) && lecturer.UserID != loggedInUserID {
		return c.Status(403).JSON(fiber.Map{"message": "Forbidden: Anda tidak boleh melihat bimbingan dosen lain"})
	}

	students, err := repoPostgre.FindLecturerAdvisees(lecturerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data mahasiswa bimbingan"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "data": students})
}

func GetLecturerService(c *fiber.Ctx) error {

	role := c.Locals("role_name")
	if role == "Dosen Wali" {
		return c.Status(403).JSON(fiber.Map{
			"message": "maaf, anda tidak bisa mengakses ini",
		})
	}

	lecturers, err := repoPostgre.GetLecturersRepo()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "gagal mengambil data dosen",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   lecturers,
	})
}

