package service

import (
	modelPostgre "sistem-prestasi/app/model/postgre"
	repository "sistem-prestasi/app/repository/postgre"
	
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)


// @Summary Get all users
// @Description Mengambil semua data user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/users [get]
func GetAllUsers(c *fiber.Ctx) error {
	users, err := repository.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengambil data user"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": users})
}



// @Summary Get user by ID
// @Description Mengambil data user berdasarkan ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /api/v1/users/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := repository.UserFindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": user})
}


// @Summary Create new user
// @Description Membuat user baru
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body modelPostgre.CreateUserRequest true "Create user request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users [post]
func CreateUser(c *fiber.Ctx) error {
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

	if err := repository.Create(newUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal membuat user (Duplicate username/email?)"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "User berhasil dibuat"})
}


// @Summary Update user
// @Description Update data user berdasarkan ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param body body modelPostgre.UpdateUserRequest true "Update user request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var req modelPostgre.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid input"})
	}

	existingUser, err := repository.UserFindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}

	if req.Username != "" { existingUser.Username = req.Username }
	if req.Email != "" { existingUser.Email = req.Email }
	if req.FullName != "" { existingUser.FullName = req.FullName }
	if req.IsActive != nil { existingUser.IsActive = *req.IsActive }

	if err := repository.Update(id, *existingUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal update user"})
	}

	if req.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		repository.UpdatePassword(id, string(hashed))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "User berhasil diupdate"})
}


// @Summary Delete user
// @Description Menghapus user berdasarkan ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := repository.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menghapus user"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "User berhasil dihapus"})
}


// @Summary Assign role to user
// @Description Mengubah role user berdasarkan ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param body body modelPostgre.AssignRoleRequest true "Assign role request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/{id}/role [put]
func AssignRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var req modelPostgre.AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid input"})
	}

	if err := repository.UpdateRole(id, req.RoleID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengubah role"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Role berhasil diubah"})
}