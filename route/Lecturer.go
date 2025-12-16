package route

import (
	"sistem-prestasi/app/service"
	"sistem-prestasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func LecturerRoute(app *fiber.App) {
	api := app.Group("/api/v1")

	lecturers := api.Group("/lecturers")
	lecturers.Use(middleware.Protect())
	lecturers.Get("/", service.GetLecturerService)
	lecturers.Get("/:id/advisees", service.GetLecturerAdvisees)
}
