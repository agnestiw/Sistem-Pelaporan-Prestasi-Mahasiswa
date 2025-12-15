package route

import (
	"sistem-prestasi/app/service"
	"sistem-prestasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func StudentRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	students := api.Group("/students")
	students.Use(middleware.Protect())
	students.Get("/", middleware.HasPermission("student:read"), service.GetAllStudentService)
	students.Get("/:id", middleware.HasPermission("student:read"), service.GetStudentByID)
	students.Get("/:id/achievements", middleware.HasPermission("student:read"), service.GetStudentAchievementDetailService)
	students.Put("/:id/advisor", middleware.HasPermission("student:update"), service.SetStudentAdvisorService)
}
