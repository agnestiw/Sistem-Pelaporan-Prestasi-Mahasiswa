package route

import (
	"sistem-prestasi/app/service"
	"sistem-prestasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func Analytics(app *fiber.App) {
	api := app.Group("/api/v1")
	reports := api.Group("/reports")

	reports.Use(middleware.Protect())
	reports.Use(middleware.HasPermission("analytics:read"))

	reports.Get("/statistics", service.GetStatisticsService)
	reports.Get("/student/:id", service.GetStudentReportService)

}
