package service

import (
	"time"

	"sistem-prestasi/helper"

	modelPostgre "sistem-prestasi/app/model/postgre"
	repoPostgre "sistem-prestasi/app/repository/postgre"
	memory "sistem-prestasi/memory"

	"github.com/gofiber/fiber/v2"
)


func Login(c *fiber.Ctx) error {
	var req modelPostgre.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	user, err := repoPostgre.Authenticate(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	permissions, _ := repoPostgre.GetPermissionsByRoleID(user.RoleID)

	accessToken, _ := helper.GenerateJWT(
		user.ID,
		user.RoleID,
		user.RoleName,
		user.StudentID, // ✅ TERISI DARI JOIN students
		permissions,
		time.Hour,
	)

	refreshToken, _ := helper.GenerateJWT(
		user.ID,
		user.RoleID,
		user.RoleName,
		user.StudentID,
		permissions,
		time.Hour*24*7,
	)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": modelPostgre.LoginResponse{
			Token:        accessToken,
			RefreshToken: refreshToken,
			User: modelPostgre.UserDetail{
				ID:          user.ID,
				Username:    user.Username,
				FullName:    user.FullName,
				RoleID:      user.RoleID,
				Role:        user.RoleName,
				Permissions: permissions,
			},
		},
	})
}


func Refresh(c *fiber.Ctx) error {
	var req modelPostgre.RefreshRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Format JSON salah. Gunakan key 'refreshToken'"})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "RefreshToken wajib diisi"})
	}

	claims, err := helper.ValidateJWT(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Refresh token tidak valid atau expired"})
	}

	userID := claims["user_id"].(string)
	roleID := claims["role_id"].(string)
	roleName := claims["role_name"].(string)
	studentID := claims["student_id"].(string)
	
	var permissions []string
	if permInter, ok := claims["permissions"].([]interface{}); ok {
		for _, p := range permInter {
			permissions = append(permissions, p.(string))
		}
	}

	newAccessToken, _ := helper.GenerateJWT(userID, roleID, roleName, &studentID, permissions, time.Hour*1)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token": newAccessToken,
		},
	})
}

func Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if len(authHeader) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Token invalid"})
	}
	
	tokenString := authHeader[7:]
	memory.AddToBlacklist(tokenString)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Berhasil logout",
	})
}

func Profile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	user, err := repoPostgre.UserFindByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}

	permissions, _ := repoPostgre.GetPermissionsByRoleID(user.RoleID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":          user.ID,
			"username":    user.Username,
			"fullName":    user.FullName,
			"email":       user.Email,
			"role":        user.RoleName,
			"permissions": permissions,
			"isActive":    user.IsActive,
			"joinedAt":    user.CreatedAt,
		},
	})
}