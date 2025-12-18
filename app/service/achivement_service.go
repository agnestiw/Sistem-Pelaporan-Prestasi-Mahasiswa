package service

import (
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	modelMongo "sistem-prestasi/app/model/mongo"
	modelPg "sistem-prestasi/app/model/postgre"
	repoMongo "sistem-prestasi/app/repository/mongo"
	repoPg "sistem-prestasi/app/repository/postgre"
)

// @Summary Get all achievements with pagination
// @Tags Achievements
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/achievements [get]
func GetAllAchievementsService(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    if page <= 0 { page = 1 }
    if limit <= 0 { limit = 10 }
    offset := (page - 1) * limit

    nama_role := c.Locals("role_name")
    var result []modelPg.AchievementReference
    var totalData int
    var err error

    if nama_role == "Mahasiswa" {
        id_mahasiswa, ok := c.Locals("student_id").(string)
        if !ok || id_mahasiswa == "" {
            return c.Status(400).JSON(fiber.Map{"message": "student_id tidak ditemukan"})
        }

        result, totalData, err = repoPg.GetAllAchievementByStudentID(id_mahasiswa, limit, offset)
    } else {
        result, totalData, err = repoPg.GetAllAchievementsRepo(limit, offset)
    }

    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "message": "Data tidak ditemukan atau error database",
            "error":   err.Error(),
        })
    }

    totalPages := int(math.Ceil(float64(totalData) / float64(limit)))

    return c.JSON(fiber.Map{
        "status": "success",
        "data":   result,
        "pagination": fiber.Map{
            "total_items":  totalData,
            "total_pages":  totalPages,
            "current_page": page,
            "limit":        limit,
        },
    })
}


// @Summary Create achievement
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body modelMongo.CreateAchievementRequest true "Create achievement"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/achievements [post]
func CreateAchievementService(c *fiber.Ctx) error {
	var req modelMongo.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	roleName := strings.ToLower(c.Locals("role_name").(string))

	var finalStudentID string

	if roleName == "admin" {
		if req.StudentID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Admin wajib menyertakan studentId",
			})
		}
		finalStudentID = req.StudentID

	} else if roleName == "mahasiswa" {
		if req.StudentID == "" {
			return c.Status(404).JSON(fiber.Map{
				"error": "Student ID Tidak valid",
			})
		}

		student, err := repoPg.GetStudentByIDRepo(req.StudentID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "Data mahasiswa tidak ditemukan",
				"message": err.Error(),
			})
		}
		finalStudentID = student.ID

	} else {
		return c.Status(403).JSON(fiber.Map{
			"error": "Role tidak diizinkan",
		})
	}

	achievement := modelMongo.Achievement{
		StudentID:       finalStudentID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		Points:          req.Points,
		Attachments:     []modelMongo.Attachment{},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoID, err := repoMongo.InsertAchievement(ctx, achievement)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save achievement"})
	}

	ref := modelPg.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          finalStudentID,
		MongoAchievementID: mongoID,
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := repoPg.CreateAchievementRef(ref); err != nil {
		_ = repoMongo.DeleteAchievement(ctx, mongoID)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save reference"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":   "Achievement draft created",
		"reference": ref,
	})
}


// @Summary Update achievement
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param achievement_id path string true "Achievement ID"
// @Param body body modelMongo.UpdateAchievementRequest true "Update achievement"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/achievements/{achievement_id} [put]
func UpdateAchievementService(c *fiber.Ctx) error {
	refID := c.Params("achievement_id")
	if refID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Achievement reference ID is required",
		})
	}

	var req modelMongo.UpdateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	userID := c.Locals("user_id").(string)
	roleName := strings.ToLower(c.Locals("role_name").(string))

	ref, err := repoPg.GetAchievementRefByID(refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Achievement not found",
		})
	}

	if ref.Status != "draft" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Achievement cannot be edited anymore",
		})
	}

	if roleName == "mahasiswa" {
		student, err := repoPg.GetStudentByIDRepo(userID)
		if err != nil || student.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{
				"error": "Forbidden",
			})
		}
	} else if roleName != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Role not allowed",
		})
	}

	update := bson.M{}

	if req.AchievementType != "" {
		update["achievementType"] = req.AchievementType
	}
	if req.Title != "" {
		update["title"] = req.Title
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.Points != nil {
		update["points"] = *req.Points
	}
	if req.Details != nil {
		update["details"] = req.Details
	}
	if req.Tags != nil {
		update["tags"] = req.Tags
	}

	update["updatedAt"] = time.Now()

	if len(update) == 1 { 
		return c.Status(400).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := repoMongo.UpdateAchievementFieldsByID(
		ctx,
		ref.MongoAchievementID,
		update,
	); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update achievement",
		})
	}

	_ = repoPg.UpdateAchievementRefUpdatedAt(refID)

	return c.JSON(fiber.Map{
		"message": "Achievement updated successfully",
	})
}


// @Summary Delete achievement
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param achievement_id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/achievements/{achievement_id} [delete]
func DeleteAchievementService(c *fiber.Ctx) error {
	refID := c.Params("achievement_id")
	if refID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Achievement reference ID is required",
		})
	}

	userID := c.Locals("user_id").(string)
	roleName := strings.ToLower(c.Locals("role_name").(string))

	ref, err := repoPg.GetAchievementRefByID(refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Achievement not found",
		})
	}

	if ref.Status == "verified" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Verified achievement cannot be deleted",
		})
	}

	if roleName == "mahasiswa" {
		student, err := repoPg.GetStudentByIDRepo(userID)
		if err != nil || student.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{
				"error": "Forbidden",
			})
		}
	} else if roleName != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Role not allowed",
		})
	}

	if err := repoPg.SoftDeleteAchievementRef(refID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete achievement",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = repoMongo.TouchAchievement(ctx, ref.MongoAchievementID)

	return c.JSON(fiber.Map{
		"message": "Achievement deleted successfully",
	})
}


// @Summary Get achievement detail
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param achievement_id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/achievements/{achievement_id} [get]
func GetAchievementDetailService(c *fiber.Ctx) error {
	achievementID := c.Params("achievement_id")
	if achievementID == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "id achievement tidak valid",
		})
	}

	result, err := repoPg.GetAchievementDetailByAchievementIDRepo(achievementID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil data achievement reference",
			"error":   err.Error(),
		})
	}

	mongoData, err := repoMongo.FindAchievementByID(
		context.Background(),
		result.MongoAchievementID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil detail achievement di MongoDB",
			"error":   err.Error(),
		})
	}

	attachments, err := repoMongo.GetAttachmentsByReferenceID(
		context.Background(),
		result.ID,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil attachments",
			"error":   err.Error(),
		})
	}

	mongoData.Attachments = attachments

	return c.Status(200).JSON(fiber.Map{
		"status":      "success",
		"reference":   result,
		"achievement": mongoData,
	})
}


// @Summary Submit achievement
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param achievement_references_id path string true "Achievement Reference ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/achievements/{achievement_references_id}/submit [post]
func SubmitAchievementService(c *fiber.Ctx) error {

	achievement_references_id := c.Params("achievement_references_id")

	student_id := c.Locals("student_id")
	roleName := c.Locals("role_name")
	fmt.Println(roleName)

	if student_id == "" {
		return c.Status(404).JSON(fiber.Map{
			"message": "id mahasiswa tidak ada",
		})
	}

	if roleName != "Admin" && roleName != "Mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Role gak boleh",
		})
	}

	studentIDfromAchievementReferences, err := repoPg.GetStudentIdFromAchievementReferences(achievement_references_id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "tidak dapat menemukan student_id di table achievement_references",
		})
	}

	if student_id != studentIDfromAchievementReferences {
		return c.Status(403).JSON(fiber.Map{
			"message": "Mahasiswa hanya boleh mengakses achievement miliknya sendiri",
		})
	}

	result, err := repoPg.SubmitAchievementRepo(achievement_references_id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak dapat submit achievement",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "berhasil submit achievement",
		"data":    result,
	})
}


// @Summary Verify achievement
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param achievement_references_id path string true "Achievement Reference ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/achievements/{achievement_references_id}/verify [post]
func VerifyAchievementService(c *fiber.Ctx) error {

	achievement_references_id := c.Params("achievement_references_id")
	roleName := c.Locals("role_name")
	
	user_id := c.Locals("user_id").(string)

	if roleName == "Mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"message": "Hanya Admin dan Dosen wali yang bisa verifikasi achievement",
		})
	}

	if roleName == "Dosen Wali" {

		advisorID, ok := c.Locals("advisor_id").(string)
		if !ok || advisorID == "" {
			return c.Status(401).JSON(fiber.Map{
				"message": "advisor_id pada token tidak valid",
			})
		}

		refAdvisorID, err := repoPg.GetAdvisorIDByAchievementRef(achievement_references_id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"message": "achievement tidak ditemukan",
				"error":   err.Error(),
			})
		}

		if advisorID != refAdvisorID {
			return c.Status(403).JSON(fiber.Map{
				"message": "Anda bukan dosen wali dari achievement ini",
			})
		}
	}

	result, err := repoPg.VerifyAchievementRepo(achievement_references_id, user_id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak dapat verify achievement_references",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "berhasil verify achievement",
		"data":    result,
	})

}


// @Summary Reject achievement
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param achievement_references_id path string true "Achievement Reference ID"
// @Param body body modelPg.RejectAchievementRequest true "Reject reason"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/achievements/{achievement_references_id}/reject [post]
func RejectAchievementService(c *fiber.Ctx) error {

	achievement_references_id := c.Params("achievement_references_id")
	roleName := c.Locals("role_name").(string)
	user_id := c.Locals("user_id").(string)

	if roleName == "Mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"message": "Hanya Admin yang bisa menolak achievement",
		})
	}

	// cek apakah dosen wali yang bersangkutan
	if roleName == "Dosen Wali" {
		advisorID, ok := c.Locals("advisor_id").(string)
		if !ok || advisorID == "" {
			return c.Status(401).JSON(fiber.Map{
				"message": "advisor_id pada token tidak valid",
			})
		}

		refAdvisorID, err := repoPg.GetAdvisorIDByAchievementRef(achievement_references_id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"message": "achievement tidak ditemukan",
				"error":   err.Error(),
			})
		}

		if advisorID != refAdvisorID {
			return c.Status(403).JSON(fiber.Map{
				"message": "Anda bukan dosen wali dari achievement ini",
			})
		}
	}

	var request modelPg.RejectAchievementRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Request body tidak valid",
		})
	}

	hasil, err := repoPg.RejectAchievementRepo(achievement_references_id, request.RejectionNote, user_id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak dapat reject achievement_references",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "berhasil reject achievement",
		"data":    hasil,
	})

}


// @Summary Upload achievement attachment
// @Tags Achievements
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param achievement_references_id path string true "Achievement Reference ID"
// @Param attachment formData file true "Attachment file"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/achievements/{achievement_references_id}/attachment [post]
func UploadAttachmentAchievementService(c *fiber.Ctx) error {
	achievementReferencesID := c.Params("achievement_references_id")
	if achievementReferencesID == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "achievement_references_id tidak boleh kosong",
		})
	}

	fileHeader, err := c.FormFile("attachment")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "File attachment tidak ditemukan",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal membuka file"})
	}
	defer file.Close()

	folder := fmt.Sprintf("./uploads/achievements/%s/", achievementReferencesID)

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, 0755)
	}

	fileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileHeader.Filename)
	filePath := folder + fileName

	dst, err := os.Create(filePath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal membuat file"})
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menyimpan file"})
	}

	folderName, err := repoMongo.UploadAttachmentAchievemenRepo(achievementReferencesID, fileName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menyimpan metadata ke database",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":   "Upload berhasil",
		"file_name": fileName,
		"folder":    folderName,
		"path":      filePath,
	})
}


// @Summary Get achievement history
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param achievement_references_id path string true "Achievement Reference ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/achievements/{achievement_references_id}/history [get]
func GetAchievementHistoryService(c *fiber.Ctx) error {
	achievement_references_id := c.Params("achievement_references_id")

	ref, err := repoPg.GetAchievementRefByID(achievement_references_id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "tidak dapat mengambil data achievement",
			"error":   err.Error(),
		})
	}
	if ref == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "terjadi kesalahan ketika mengambil data achievement",
			"error":   err.Error(),
		})
	}

	achievement, err := repoMongo.FindAchievementByID(context.Background(), ref.MongoAchievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "tidak dapat mengambil data achievement",
			"error":   err.Error(),
		})
	}

	history := []modelPg.HistoryItem{}

	history = append(history, modelPg.HistoryItem{
		Status:    "draft",
		Timestamp: ref.CreatedAt,
		Note:      "",
	})

	if ref.SubmittedAt != nil {
		history = append(history, modelPg.HistoryItem{
			Status:    "submitted",
			Timestamp: *ref.SubmittedAt,
			Note:      "",
		})
	}

	if ref.VerifiedAt != nil {
		history = append(history, modelPg.HistoryItem{
			Status:    "verified",
			Timestamp: *ref.VerifiedAt,
			Note:      "",
		})
	}

	if ref.Status == "rejected" {
		history = append(history, modelPg.HistoryItem{
			Status:    "rejected",
			Timestamp: ref.UpdatedAt,
			Note: func() string { // FIX #1
				if ref.RejectionNote != nil {
					return *ref.RejectionNote
				}
				return ""
			}(),
		})
	}

	response := &modelPg.HistoryResponse{
		Reference:   ref,
		Achievement: &achievement,
		History:     history,
	}

	return c.Status(200).JSON(fiber.Map{
		"data":   response,
		"status": "success",
	})
}
