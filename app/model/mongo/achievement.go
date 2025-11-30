package mongo

import "time"

type Achievement struct {
	ID              string                 `bson:"_id,omitempty" json:"id"`
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