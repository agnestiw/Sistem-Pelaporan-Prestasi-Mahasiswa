package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	StudentID       string                 `bson:"studentId" json:"studentId"`
	AchievementType string                 `bson:"achievementType" json:"achievementType"`
	Title           string                 `bson:"title" json:"title"`
	Description     string                 `bson:"description" json:"description"`
	Details         map[string]interface{} `bson:"details" json:"details"` 
	Attachments     []Attachment           `bson:"attachments" json:"attachments"`
	Tags            []string               `bson:"tags" json:"tags"`
	Points          int                    `bson:"points" json:"points"`
	CreatedAt       time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time              `bson:"updatedAt" json:"updatedAt"`
}

type AchievementResponse struct {
	ID                 string                 `json:"id"`                   
	MongoID            string                 `json:"mongo_achievement_id"`
	StudentID          string                 `json:"student_id"`
	StudentName        string                 `json:"student_name"`
	Status             string                 `json:"status"`
	AchievementType    string                 `json:"achievement_type"`
	Details            map[string]interface{} `json:"details,omitempty"`
	Points             int                    `json:"points"`
	Tags 			   []string 			  `json:"tags,omitempty"`
	CreatedAt          time.Time              `json:"created_at"`
	SubmittedAt        *time.Time             `json:"submitted_at"`
	VerifiedAt         *time.Time             `json:"verified_at"`
	VerifiedBy         *string                `json:"verified_by"`
	RejectionNote      *string                `json:"rejection_note"`
}

type CreateAchievementRequest struct {
	StudentID       string                 `json:"studentId"`
	AchievementType string                 `json:"achievementType"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Points          int                    `json:"points"`
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags"`
}

type UpdateAchievementRequest struct {
	AchievementType string                 `json:"achievementType,omitempty"`
	Title           string                 `json:"title,omitempty"`
	Description     string                 `json:"description,omitempty"`
	Points          *int                   `json:"points,omitempty"`
	Details         map[string]interface{} `json:"details,omitempty"`
	Tags            []string               `json:"tags,omitempty"`
}
