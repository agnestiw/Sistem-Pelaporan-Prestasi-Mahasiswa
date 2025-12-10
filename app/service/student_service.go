package service

import (
	modelPostgre "sistem-prestasi/app/model/postgre"
	repoPostgre "sistem-prestasi/app/repository/postgre"
	"sistem-prestasi/helper"

	"github.com/gofiber/fiber/v2"
)

// GetMyStudents: Digunakan oleh Dosen untuk melihat mahasiswa bimbingannya sendiri
// func GetMyStudents(c *fiber.Ctx) error {


    // loggedInUserID := c.Locals("user_id").(string)
	
	

   
    // lecturer, err := repoPostgre.FindLecturerByUserID(loggedInUserID)
    // if err != nil {
    //     return c.Status(404).JSON(fiber.Map{
    //         "message": "Data dosen tidak ditemukan untuk user ini. Apakah anda login sebagai Dosen?",
    //     })
    // }

    // // 3. Ambil mahasiswa bimbingan menggunakan ID Dosen yang ditemukan
    // students, err := repoPostgre.FindLecturerAdvisees(lecturer.ID)
    // if err != nil {
    //     return c.Status(500).JSON(fiber.Map{
    //         "message": "Gagal mengambil data mahasiswa bimbingan",
    //         "error":   err.Error(),
    //     })
    // }

    // return c.Status(200).JSON(fiber.Map{"status": "success", "data": students})
// }




func GetAll(c *fiber.Ctx) error {

	// ngambil role yang di jwt 
	nama_role := c.Locals("role_name")
	if nama_role == "Mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"message": "anda bukan seorang admin maupun dosen",
		})
	}
	
	if nama_role == "Dosen Wali" {
		
	}


	students, err := repoPostgre.StudentFindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil data mahasiswa",
			"error":   err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "data": students})
}


func StudentFindByID(c *fiber.Ctx) error {
	id := c.Params("id")

	student, err := repoPostgre.StudentFindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Mahasiswa tidak ditemukan"})
	}

	loggedInUserID := c.Locals("user_id").(string)

	if !helper.IsAdmin(c) && student.UserID != loggedInUserID {
		return c.Status(403).JSON(fiber.Map{"message": "Forbidden: Anda tidak boleh melihat profil mahasiswa lain"})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "data": student})
}

func AssignAdvisor(c *fiber.Ctx) error {
	id := c.Params("id")

	var req modelPostgre.AssignAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid input"})
	}

	student, err := repoPostgre.StudentFindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Mahasiswa tidak ditemukan"})
	}

	loggedInUserID := c.Locals("user_id").(string)
	if !helper.IsAdmin(c) && student.UserID != loggedInUserID {
		return c.Status(403).JSON(fiber.Map{"message": "Forbidden: Anda hanya boleh memilih dosen wali untuk diri sendiri"})
	}

	if err := repoPostgre.UpdateAdvisor(id, req.LecturerID); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal set dosen wali"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Dosen Wali berhasil diupdate"})
}