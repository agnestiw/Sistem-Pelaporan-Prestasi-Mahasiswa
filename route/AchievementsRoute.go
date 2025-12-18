package route

import (
	"sistem-prestasi/app/service"
	"sistem-prestasi/middleware"

	"github.com/gofiber/fiber/v2"
)

func AchievementRoutes(app *fiber.App) {

	api := app.Group("/api/v1")

	achievments := api.Group("/achievements")
	achievments.Use(middleware.Protect())
	achievments.Get("/", middleware.HasPermission("achievement:read"), service.GetAllAchievementsService)
	achievments.Get("/:achievement_id", middleware.HasPermission("achievement:read"), service.GetAchievementDetailService)
	achievments.Post("/", middleware.HasPermission("achievement:create"), service.CreateAchievementService)
	achievments.Put("/:achievement_id", middleware.HasPermission("achievement:update"), service.UpdateAchievementService)
	achievments.Delete("/:achievement_id", middleware.HasPermission("achievement:delete"), service.DeleteAchievementService)

	achievments.Post("/:achievement_references_id/submit", middleware.HasPermission("achievement:submit"), service.SubmitAchievementService)
	achievments.Post("/:achievement_references_id/verify", middleware.HasPermission("achievement:verify"), service.VerifyAchievementService)
	achievments.Post("/:achievement_references_id/reject", middleware.HasPermission("achievement:reject"), service.RejectAchievementService)
	achievments.Post("/:achievement_references_id/attachment", middleware.HasPermission("achievement:create"), service.UploadAttachmentAchievementService)
	achievments.Get("/:achievement_references_id/history", middleware.HasPermission("achievement:read"), service.GetAchievementHistoryService)

}
