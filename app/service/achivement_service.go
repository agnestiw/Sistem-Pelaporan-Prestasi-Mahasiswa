package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	modelMongo "sistem-prestasi/app/model/mongo"
	modelPg "sistem-prestasi/app/model/postgre"
	repoMongo "sistem-prestasi/app/repository/mongo"
	repoPg "sistem-prestasi/app/repository/postgre"
)

func GetAllAchievementsService(c *fiber.Ctx) error {

	// kalo role kosong
	nama_role := c.Locals("role_name")

	// jika mahasiswa
	if nama_role == "Mahasiswa" {
		id_mahasiswa := c.Locals("student_id").(string)

		result, err := repoPg.GetAllAchievementByStudentID(id_mahasiswa)

		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"message": "tidak bisa boss",
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"data":   result,
		})
	}

	result, err := repoPg.GetAllAchievementsRepo()
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "tidak bisa boss",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func CreateAchievementService(c *fiber.Ctx) error {
	var input modelMongo.Achievement
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	roleName, ok := c.Locals("role_name").(string)
	if !ok || roleName == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	roleName = strings.ToLower(roleName)

	var finalStudentID string

	if roleName == "admin" {
		if input.StudentID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Admin wajib menyertakan studentId",
			})
		}
		finalStudentID = input.StudentID

	} else if roleName == "mahasiswa" {
		if input.StudentID != "" {
			return c.Status(403).JSON(fiber.Map{
				"error": "Mahasiswa tidak boleh mengirim studentId",
			})
		}

		studentID, err := repoPg.GetStudentByUserID(userID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "Data mahasiswa tidak ditemukan",
			})
		}

		finalStudentID = studentID
	} else {
		return c.Status(403).JSON(fiber.Map{
			"error": "Role tidak diizinkan membuat achievement",
		})
	}

	input.StudentID = finalStudentID
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoID, err := repoMongo.InsertAchievement(ctx, input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save to Mongo"})
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
		_ = repoMongo.DeleteAchievement(ctx, mongoID) // rollback
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save reference"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Achievement draft created",
		"data":    ref,
	})
}

// func UpdateAchievementService(c *fiber.Ctx) error {
// 	refID := c.Params("id")
// 	if refID == "" {
// 		return c.Status(400).JSON(fiber.Map{
// 			"error": "Achievement reference ID is required",
// 		})
// 	}

// 	var input modelMongo.Achievement
// 	if err := c.BodyParser(&input); err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"error": "Invalid input",
// 		})
// 	}

// 	userID, ok := c.Locals("user_id").(string)
// 	if !ok || userID == "" {
// 		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
// 	}

// 	roleName, ok := c.Locals("role_name").(string)
// 	if !ok || roleName == "" {
// 		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
// 	}

// 	roleName = strings.ToLower(roleName)

// 	// 🔹 Ambil reference dari PostgreSQL
// 	ref, err := repoPg.GetAchievementRefByID(refID)
// 	if err != nil {
// 		return c.Status(404).JSON(fiber.Map{
// 			"error": "Achievement reference not found",
// 		})
// 	}

// 	// 🔒 Authorization check
// 	if roleName == "mahasiswa" {
// 		studentID, err := repoPg.GetStudentByUserID(userID)
// 		if err != nil || studentID != ref.StudentID {
// 			return c.Status(403).JSON(fiber.Map{
// 				"error": "Forbidden",
// 			})
// 		}
// 	} else if roleName != "admin" {
// 		return c.Status(403).JSON(fiber.Map{
// 			"error": "Role not allowed",
// 		})
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// 🔹 Update Mongo Data
// 	input.UpdatedAt = time.Now()
// 	err = mongoRepo.UpdateByID(ctx, ref.MongoAchievementID, input)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"error": "Failed to update mongo achievement",
// 		})
// 	}

// 	// 🔹 Update reference updated_at
// 	if err := repoPg.UpdateAchievementRefUpdatedAt(refID); err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"error": "Failed to update reference",
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"message": "Achievement updated successfully",
// 	})
// }

func GetAchievementDetailService(c *fiber.Ctx) error {
	id := c.Params("id")

	ref, err := repoPg.GetAchievementRefByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Achievement not found",
		})
	}

	ctx := context.Background()
	achievement, err := repoMongo.FindAchievementByID(ctx, ref.MongoAchievementID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":    "Detail data missing",
			"mongo_ID": ref.MongoAchievementID,
			"message":  err.Error(),
		})
	}

	response := modelMongo.AchievementResponse{
		ID:          ref.ID,
		MongoID:     ref.MongoAchievementID,
		StudentID:   ref.StudentID,
		StudentName: ref.StudentName,
		Status:      ref.Status,
		Details:     achievement.Details,
	}

	return c.JSON(fiber.Map{"data": response})
}

func SubmitAchievementService(c *fiber.Ctx) error {

	achievement_references_id := c.Params("achievement_references_id")

	student_id := c.Locals("student_id")
	roleName := c.Locals("role_name")
	fmt.Println(roleName)

	// kalo id mahasiswa gak ada
	if student_id == "" {
		return c.Status(404).JSON(fiber.Map{
			"message": "id mahasiswa tidak ada",
		})
	}

	// cek kalau dosen atau tidak
	if roleName != "Admin" && roleName != "Mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Role gak boleh",
		})
	}

	// cek apakah student_id yang login sesuai dengan student_id di table achievement_references
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

func VerifyAchievementService(c *fiber.Ctx) error {

	achievement_references_id := c.Params("achievement_references_id")
	roleName := c.Locals("role_name")

	if roleName == "Mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"message": "Hanya Admin yang bisa verifikasi achievement",
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

	result, err := repoPg.VerifyAchievementRepo(achievement_references_id)
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

func UploadAttachmentAchievementService(c *fiber.Ctx) error {
	achievementReferencesID := c.FormValue("achievement_references_id")
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

	// Folder: /uploads/achievements/<achievement_references_id>/
	folder := fmt.Sprintf("./uploads/achievements/%s/", achievementReferencesID)

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, 0755)
	}

	// file name unique
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

	// Save metadata
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

	// 3. generate history
	history := []modelPg.HistoryItem{}

	// DRAFT
	history = append(history, modelPg.HistoryItem{
		Status:    "draft",
		Timestamp: ref.CreatedAt,
		Note:      "",
	})

	// SUBMITTED
	if ref.SubmittedAt != nil {
		history = append(history, modelPg.HistoryItem{
			Status:    "submitted",
			Timestamp: *ref.SubmittedAt,
			Note:      "",
		})
	}

	// VERIFIED
	if ref.VerifiedAt != nil {
		history = append(history, modelPg.HistoryItem{
			Status:    "verified",
			Timestamp: *ref.VerifiedAt,
			Note:      "",
		})
	}

	// REJECTED
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
