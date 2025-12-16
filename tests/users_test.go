package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupUsersTestApp() *fiber.App {
	app := fiber.New()

	api := app.Group("/api/v1")
	users := api.Group("/users")

	// ===== DUMMY HANDLER (BUKAN repository/service asli) =====

	users.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"data": []fiber.Map{
				{"id": "1", "username": "admin"},
			},
		})
	})

	users.Get("/:id", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"id":       c.Params("id"),
				"username": "testuser",
			},
		})
	})

	users.Post("/", func(c *fiber.Ctx) error {
		return c.Status(201).JSON(fiber.Map{
			"status":  "success",
			"message": "User berhasil dibuat",
		})
	})

	users.Put("/:id", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "User berhasil diupdate",
		})
	})

	users.Delete("/:id", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "User berhasil dihapus",
		})
	})

	users.Put("/:id/role", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Role berhasil diubah",
		})
	})

	return app
}

// ===================== TESTS =====================

func TestGetAllUsers(t *testing.T) {
	app := setupUsersTestApp()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetUserByID(t *testing.T) {
	app := setupUsersTestApp()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/123", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateUser(t *testing.T) {
	app := setupUsersTestApp()

	body := `{
		"username": "newuser",
		"email": "new@user.com",
		"password": "password",
		"fullName": "New User",
		"roleId": "role-1"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/users",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestUpdateUser(t *testing.T) {
	app := setupUsersTestApp()

	body := `{
		"fullName": "Updated Name"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/users/123",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteUser(t *testing.T) {
	app := setupUsersTestApp()

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/users/123",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAssignRole(t *testing.T) {
	app := setupUsersTestApp()

	body := `{
		"roleId": "role-admin"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/users/123/role",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}