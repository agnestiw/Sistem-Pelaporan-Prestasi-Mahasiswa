package service

import (
	modelPostgre "sistem-prestasi/app/model/postgre"
	repoPostgre "sistem-prestasi/app/repository/postgre"
	"github.com/gofiber/fiber/v2"
)

type StudentService struct {
	StudentRepo *repoPostgre.StudentRepository
}

func NewStudentService(r *repoPostgre.StudentRepository) *StudentService {
	return &StudentService{StudentRepo: r}
}

func (s *StudentService) GetAll(c *fiber.Ctx) error {
	students, err := s.StudentRepo.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data mahasiswa"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "data": students})
}

func (s *StudentService) isAdmin(c *fiber.Ctx) bool {
	perms, ok := c.Locals("permissions").([]string)
	if !ok { return false }
	for _, p := range perms {
		if p == "user:manage" { return true }
	}
	return false
}

func (s *StudentService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	student, err := s.StudentRepo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Mahasiswa tidak ditemukan"})
	}

	loggedInUserID := c.Locals("user_id").(string)

	if !s.isAdmin(c) && student.UserID != loggedInUserID {
		return c.Status(403).JSON(fiber.Map{"message": "Forbidden: Anda tidak boleh melihat profil mahasiswa lain"})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "data": student})
}

func (s *StudentService) AssignAdvisor(c *fiber.Ctx) error {
	id := c.Params("id")

	var req modelPostgre.AssignAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid input"})
	}

	student, err := s.StudentRepo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Mahasiswa tidak ditemukan"})
	}

	loggedInUserID := c.Locals("user_id").(string)
	if !s.isAdmin(c) && student.UserID != loggedInUserID {
		return c.Status(403).JSON(fiber.Map{"message": "Forbidden: Anda hanya boleh memilih dosen wali untuk diri sendiri"})
	}

	if err := s.StudentRepo.UpdateAdvisor(id, req.LecturerID); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal set dosen wali"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Dosen Wali berhasil diupdate"})
}