package service

import (
	repoPostgre "sistem-prestasi/app/repository/postgre"
	"github.com/gofiber/fiber/v2"
)

type LecturerService struct {
	LecturerRepo *repoPostgre.LecturerRepository
}

func NewLecturerService(r *repoPostgre.LecturerRepository) *LecturerService {
	return &LecturerService{LecturerRepo: r}
}

func (s *LecturerService) GetAll(c *fiber.Ctx) error {
	lecturers, err := s.LecturerRepo.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data dosen"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "data": lecturers})
}

func (s *LecturerService) isAdmin(c *fiber.Ctx) bool {
	perms, ok := c.Locals("permissions").([]string)
	if !ok { return false }
	for _, p := range perms {
		if p == "user:manage" { return true }
	}
	return false
}

func (s *LecturerService) GetAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")

	lecturer, err := s.LecturerRepo.FindByID(lecturerID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Dosen tidak ditemukan"})
	}

	loggedInUserID := c.Locals("user_id").(string)

	if !s.isAdmin(c) && lecturer.UserID != loggedInUserID {
		return c.Status(403).JSON(fiber.Map{"message": "Forbidden: Anda tidak boleh melihat bimbingan dosen lain"})
	}

	students, err := s.LecturerRepo.FindAdvisees(lecturerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data mahasiswa bimbingan"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "data": students})
}