package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupAchievementsTestApp() *fiber.App {
	app := fiber.New()

	api := app.Group("/api/v1")
	achievements := api.Group("/achievements")

	// ===== DUMMY HANDLERS =====

	achievements.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"data":   []fiber.Map{},
		})
	})

	achievements.Get("/:achievement_id", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"id":     c.Params("achievement_id"),
		})
	})

	achievements.Post("/", func(c *fiber.Ctx) error {
		return c.Status(201).JSON(fiber.Map{
			"message": "Achievement draft created",
		})
	})

	achievements.Put("/:achievement_id", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "Achievement updated successfully",
		})
	})

	achievements.Delete("/:achievement_id", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "Achievement deleted successfully",
		})
	})

	achievements.Post("/:achievement_references_id/submit", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "berhasil submit achievement",
		})
	})

	achievements.Post("/:achievement_references_id/verify", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "berhasil verify achievement",
		})
	})

	achievements.Post("/:achievement_references_id/reject", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "berhasil reject achievement",
		})
	})

	achievements.Post("/:achievement_references_id/attachment", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "Upload berhasil",
		})
	})

	achievements.Get("/:achievement_references_id/history", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"data":   []fiber.Map{},
		})
	})

	return app
}

// ===================== TESTS =====================

func TestGetAllAchievements(t *testing.T) {
	app := setupAchievementsTestApp()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/achievements", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetAchievementDetail(t *testing.T) {
	app := setupAchievementsTestApp()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/achievements/123", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateAchievement(t *testing.T) {
	app := setupAchievementsTestApp()

	body := `{
		"title": "Juara 1",
		"achievementType": "competition"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/achievements",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestUpdateAchievement(t *testing.T) {
	app := setupAchievementsTestApp()

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/achievements/123",
		strings.NewReader(`{"title":"Updated"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteAchievement(t *testing.T) {
	app := setupAchievementsTestApp()

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/achievements/123",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSubmitAchievement(t *testing.T) {
	app := setupAchievementsTestApp()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/achievements/123/submit",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestVerifyAchievement(t *testing.T) {
	app := setupAchievementsTestApp()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/achievements/123/verify",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRejectAchievement(t *testing.T) {
	app := setupAchievementsTestApp()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/achievements/123/reject",
		strings.NewReader(`{"reason":"Tidak valid"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUploadAttachmentAchievement(t *testing.T) {
	app := setupAchievementsTestApp()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/achievements/123/attachment",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetAchievementHistory(t *testing.T) {
	app := setupAchievementsTestApp()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/achievements/123/history",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
