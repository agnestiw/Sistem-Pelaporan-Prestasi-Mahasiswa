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

	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)

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
}