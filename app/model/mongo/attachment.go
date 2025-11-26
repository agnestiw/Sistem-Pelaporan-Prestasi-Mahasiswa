package mongo

import "time"

type Attachment struct {
	FileName   string    `bson:"fileName" json:"fileName"`
	FileURL    string    `bson:"fileUrl" json:"fileUrl"`
	FileType   string    `bson:"fileType" json:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploadedAt"`
}