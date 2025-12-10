package service

import (
	"context"
	"fmt"
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
