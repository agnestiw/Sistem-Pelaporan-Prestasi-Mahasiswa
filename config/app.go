package config

import (
	"sistem-prestasi/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "sistem-prestasi/docs"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New(LoggerConfig()))

	// swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// routes
	route.UsersRoute(app)
	route.AuthRoute(app)
	route.AchievementRoutes(app)
	route.LecturerRoute(app)
	route.StudentRoutes(app)
	route.Analytics(app)

	return app
}