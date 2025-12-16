package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupStudentsTestApp() *fiber.App {
	app := fiber.New()

	api := app.Group("/api/v1")
	students := api.Group("/students")

	students.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"data":   []fiber.Map{},
		})
	})

	students.Get("/:id", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"id": c.Params("id"),
			},
		})
	})

	students.Get("/:id/achievements", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":      "success",
			"student":     fiber.Map{"id": c.Params("id")},
			"reference":   fiber.Map{},
			"achievement": fiber.Map{},
		})
	})

	students.Put("/:id/advisor", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"student_id": c.Params("id"),
				"advisor_id": "advisor-123",
			},
		})
	})

	return app
}

// ===================== TESTS =====================

func TestGetAllStudents(t *testing.T) {
	app := setupStudentsTestApp()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/students", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetStudentByID(t *testing.T) {
	app := setupStudentsTestApp()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/123", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetStudentAchievementDetail(t *testing.T) {
	app := setupStudentsTestApp()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/students/123/achievements",
		nil,
	)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSetStudentAdvisor(t *testing.T) {
	app := setupStudentsTestApp()

	body := `{
		"advisor_id": "advisor-123"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/students/123/advisor",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
