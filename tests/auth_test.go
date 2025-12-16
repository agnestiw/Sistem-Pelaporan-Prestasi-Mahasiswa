package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	app := fiber.New()

	api := app.Group("/api/v1")
	auth := api.Group("/auth")

	auth.Post("/login", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"token": "dummy-access-token",
			},
		})
	})

	auth.Post("/refresh", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"token": "dummy-new-access-token",
			},
		})
	})

	auth.Get("/profile", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"id":       "user-123",
				"username": "testuser",
				"role":     "Mahasiswa",
			},
		})
	})

	auth.Post("/logout", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Berhasil logout",
		})
	})

	return app
}


func TestLogin(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/login",
		nil,
	)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRefresh(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/refresh",
		nil,
	)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestProfile(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/auth/profile",
		nil,
	)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestLogout(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/logout",
		nil,
	)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
