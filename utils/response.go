package utils

import "github.com/gofiber/fiber/v2"

// APIResponse adalah format standar output JSON
type APIResponse struct {
	Status  string      `json:"status"`  // "success" atau "error"
	Message string      `json:"message"` // Pesan deskriptif
	Data    interface{} `json:"data,omitempty"` // Data payload (bisa kosong jika error)
}

// SuccessResponse helper untuk return 200/201
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// ErrorResponse helper untuk return 400/401/403/500
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(APIResponse{
		Status:  "error",
		Message: message,
		Data:    nil,
	})
}