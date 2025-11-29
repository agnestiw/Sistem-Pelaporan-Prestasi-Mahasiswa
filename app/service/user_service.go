package service

import (
	modelPostgre "sistem-prestasi/app/model/postgre"
	repoPostgre "sistem-prestasi/app/repository/postgre"
	
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repoPostgre.UserRepository
}

func NewUserService(userRepo *repoPostgre.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) GetAllUsers(c *fiber.Ctx) error {
	users, err := s.UserRepo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengambil data user"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": users})
}

func (s *UserService) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := s.UserRepo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": user})
}

func (s *UserService) CreateUser(c *fiber.Ctx) error {
	var req modelPostgre.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid input"})
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	newUser := modelPostgre.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
		FullName:     req.FullName,
		RoleID:       req.RoleID,
	}

	if err := s.UserRepo.Create(newUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal membuat user (Duplicate username/email?)"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "User berhasil dibuat"})
}

func (s *UserService) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var req modelPostgre.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid input"})
	}

	existingUser, err := s.UserRepo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}

	if req.Username != "" { existingUser.Username = req.Username }
	if req.Email != "" { existingUser.Email = req.Email }
	if req.FullName != "" { existingUser.FullName = req.FullName }
	if req.IsActive != nil { existingUser.IsActive = *req.IsActive }

	if err := s.UserRepo.Update(id, *existingUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal update user"})
	}

	if req.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		s.UserRepo.UpdatePassword(id, string(hashed))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "User berhasil diupdate"})
}

func (s *UserService) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := s.UserRepo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menghapus user"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "User berhasil dihapus"})
}

func (s *UserService) AssignRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var req modelPostgre.AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid input"})
	}

	if err := s.UserRepo.UpdateRole(id, req.RoleID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengubah role"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Role berhasil diubah"})
}