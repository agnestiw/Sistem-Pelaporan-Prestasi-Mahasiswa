package route

import (
	"sistem-prestasi/app/service"
	repoPostgre "sistem-prestasi/app/repository/postgre"
	"sistem-prestasi/database"
	"sistem-prestasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	userRepo := repoPostgre.NewUserRepository(database.DB)
	studentRepo := repoPostgre.NewStudentRepository(database.DB)
	lecturerRepo := repoPostgre.NewLecturerRepository(database.DB)

	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	studentService := service.NewStudentService(studentRepo)    
	lecturerService := service.NewLecturerService(lecturerRepo)

	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/login", authService.Login)
	auth.Post("/refresh", authService.Refresh)
	auth.Post("/logout", middleware.Protect(), authService.Logout)
	auth.Get("/profile", middleware.Protect(), authService.Profile)

	users := api.Group("/users")
	users.Use(middleware.Protect())
	users.Use(middleware.HasPermission("user:manage"))
	users.Get("/", userService.GetAllUsers)
	users.Post("/", userService.CreateUser)
	users.Get("/:id", userService.GetUserByID)
	users.Put("/:id", userService.UpdateUser)
	users.Delete("/:id", userService.DeleteUser)
	users.Put("/:id/role", userService.AssignRole)

	students := api.Group("/students")
	students.Use(middleware.Protect()) 
	students.Get("/", middleware.HasPermission("user:manage"), studentService.GetAll)
	students.Get("/:id", studentService.GetByID)
	students.Put("/:id/advisor", studentService.AssignAdvisor) 
	// students.Get("/:id/achievements"), studentService.GetAchievements)

	lecturers := api.Group("/lecturers")
	lecturers.Use(middleware.Protect())
	lecturers.Get("/", middleware.HasPermission("user:manage"), lecturerService.GetAll)
	lecturers.Get("/:id/advisees", lecturerService.GetAdvisees)
}