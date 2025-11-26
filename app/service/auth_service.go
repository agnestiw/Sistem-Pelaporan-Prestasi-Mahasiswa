package service

import (
	"os"
	"time"

	"sistem-prestasi/app/dto"
	"sistem-prestasi/app/repository/postgre" 

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *postgre.UserRepository
}

func NewAuthService(userRepo *postgre.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	user, err := s.UserRepo.FindByUsername(req.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Username atau password salah",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Username atau password salah",
		})
	}

	permissions, err := s.UserRepo.GetPermissionsByRoleID(user.RoleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data permission",
		})
	}

	token, err := s.generateJWT(user.ID, user.RoleID, permissions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login berhasil",
		"data": fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":          user.ID,
				"username":    user.Username,
				"fullName":    user.FullName,
				"role":        user.RoleName,
				"permissions": permissions, 
			},
		},
	})
}

func (s *AuthService) generateJWT(userID, roleID string, permissions []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":     userID,
		"role_id":     roleID,
		"permissions": permissions,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("API_SECRET")
	return token.SignedString([]byte(secret))
}
