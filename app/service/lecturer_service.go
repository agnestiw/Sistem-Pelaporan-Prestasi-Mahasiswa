package service

import (
	repoPostgre "sistem-prestasi/app/repository/postgre"
	"sistem-prestasi/helper"

	"github.com/gofiber/fiber/v2"
)

// GetAllLecturers mengambil semua data dosen
func GetAllLecturers(c *fiber.Ctx) error {
	lecturers, err := repoPostgre.FindAllLecturers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data dosen"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "data": lecturers})
}

// GetLecturerAdvisees mengambil data mahasiswa bimbingan dosen tertentu
func GetLecturerAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")

	// Menggunakan fungsi repo functional: FindLecturerByID(db, id)
	lecturer, err := repoPostgre.FindLecturerByID(lecturerID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Dosen tidak ditemukan"})
	}

	loggedInUserID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
	}

	// Cek otorisasi: harus admin atau dosen yang bersangkutan
	if !helper.IsAdmin(c) && lecturer.UserID != loggedInUserID {
		return c.Status(403).JSON(fiber.Map{"message": "Forbidden: Anda tidak boleh melihat bimbingan dosen lain"})
	}

	// Menggunakan fungsi repo functional: FindLecturerAdvisees(db, lecturerID)
	students, err := repoPostgre.FindLecturerAdvisees(lecturerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data mahasiswa bimbingan"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "data": students})
}