package service

import (
	repoPostgre "sistem-prestasi/app/repository/postgre"
	"sistem-prestasi/helper"

	"github.com/gofiber/fiber/v2"
)

// GetMyAdvisor: Digunakan oleh Mahasiswa untuk melihat siapa dosen walinya
func GetMyAdvisor(c *fiber.Ctx) error {
    // 1. Ambil User ID dari Token
    loggedInUserID := c.Locals("user_id").(string)

    // 2. Cari data Mahasiswa berdasarkan User ID
    student, err := repoPostgre.FindStudentByUserID(loggedInUserID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "message": "Data mahasiswa tidak ditemukan. Apakah anda login sebagai Mahasiswa?",
        })
    }

    // 3. Cek apakah sudah punya dosen wali
    if student.AdvisorName == nil {
         return c.Status(200).JSON(fiber.Map{
            "status": "success", 
            "message": "Anda belum memiliki Dosen Wali",
            "data": nil,
        })
    }

    // Jika ingin mengembalikan detail dosen lengkap, anda bisa query lagi ke tabel lecturers
    // Tapi jika cukup nama saja (dari query student detail), kembalikan ini:
    return c.Status(200).JSON(fiber.Map{
        "status": "success", 
        "data": fiber.Map{
            "advisorName": *student.AdvisorName,
            // Anda bisa menambahkan info lain jika query FindStudentByUserID dimodifikasi
        },
    })
}



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