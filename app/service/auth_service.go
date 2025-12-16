package service

import (
	"time"

	"sistem-prestasi/helper"

	modelPostgre "sistem-prestasi/app/model/postgre"
	repoPostgre "sistem-prestasi/app/repository/postgre"
	memory "sistem-prestasi/memory"

	"github.com/gofiber/fiber/v2"
)


// @Summary Login user
// @Description Login menggunakan username dan password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body modelPostgre.LoginRequest true "Login request"
// @Success 200 {object} modelPostgre.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/login [post]
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
		user.AdvisorID,
		permissions,
		time.Hour,
	)

	refreshToken, _ := helper.GenerateJWT(
		user.ID,
		user.RoleID,
		user.RoleName,
		user.StudentID,
		user.AdvisorID,
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



// @Summary Refresh access token
// @Description Generate access token baru menggunakan refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body modelPostgre.RefreshRequest true "Refresh token request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/refresh [post]
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
	advisor_id := claims["advisor_id"].(string)
	
	var permissions []string
	if permInter, ok := claims["permissions"].([]interface{}); ok {
		for _, p := range permInter {
			permissions = append(permissions, p.(string))
		}
	}

	newAccessToken, _ := helper.GenerateJWT(userID, roleID, roleName, &studentID, &advisor_id, permissions, time.Hour*1)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token": newAccessToken,
		},
	})
}


// @Summary Logout user
// @Description Logout user dan blacklist access token
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/auth/logout [post]
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


// @Summary Get user profile
// @Description Mengambil data profile user yang sedang login
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/auth/profile [get]
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