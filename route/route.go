package route

import (
	service "sistem-prestasi/app/service"
	"sistem-prestasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/login", service.Login)
	auth.Post("/refresh", service.Refresh)
	auth.Post("/logout", middleware.Protect(), service.Logout)
	auth.Get("/profile", middleware.Protect(), service.Profile)

	users := api.Group("/users")
	users.Use(middleware.Protect())
	users.Use(middleware.HasPermission("user:manage"))
	users.Get("/", service.GetAllUsers)
	users.Post("/", service.CreateUser)
	users.Get("/:id", service.GetUserByID)
	users.Put("/:id", service.UpdateUser)
	users.Delete("/:id", service.DeleteUser)
	users.Put("/:id/role", service.AssignRole)

	achievments := api.Group("/achievements")
	achievments.Use(middleware.Protect())
	achievments.Get("/", service.GetAllAchievementsService)
	achievments.Post("/", service.CreateAchievementService)
	achievments.Get("/:id", service.GetAchievementDetailService)
	achievments.Post("/:achievement_references_id/submit", service.SubmitAchievementService)
	achievments.Post("/:achievement_references_id/verify", service.VerifyAchievementService)

	students := api.Group("/students")
	students.Use(middleware.Protect())
	students.Get("/", service.GetAll) 
	students.Get("/:id", service.StudentFindByID)
	students.Put("/:id/advisor", service.AssignAdvisor)

	lecturers := api.Group("/lecturers")
	lecturers.Use(middleware.Protect())
	lecturers.Get("/", service.GetMyAdvisor) 
	lecturers.Get("/:id/advisees", service.GetLecturerAdvisees)
}
