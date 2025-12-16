package mongo

import (
	"context"
	"fmt"
	"path/filepath"
	model "sistem-prestasi/app/model/mongo"
	"sistem-prestasi/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepo struct {
	Collection *mongo.Collection
}

func NewAchievementRepo(db *mongo.Database) *AchievementRepo {
	return &AchievementRepo{
		Collection: db.Collection("achievements"),
	}
}

func (r *AchievementRepo) Insert(ctx context.Context, data model.Achievement) (string, error) {
	result, err := r.Collection.InsertOne(ctx, data)
	if err != nil {
		return "", err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}

func InsertAchievement(ctx context.Context, input model.Achievement) (string, error) {
	collection := database.MongoDb.Collection("achievements")

	result, err := collection.InsertOne(ctx, input)
	if err != nil {
		return "", err
	}

	oid := result.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}

func DeleteAchievement(ctx context.Context, mongoID string) error {
	collection := database.MongoDb.Collection("achievements")

	oid, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func TouchAchievement(ctx context.Context, mongoID string) error {
	collection := database.MongoDb.Collection("achievements")

	oid, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return err
	}

	_, err = collection.UpdateByID(ctx, oid, bson.M{
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	})

	return err
}

func UploadAttachmentAchievemenRepo(achievementReferencesID string, fileName string) (string, error) {
    collection := database.MongoDb.Collection("achievement_attachments")

    folder := fmt.Sprintf("%s", achievementReferencesID)

    doc := bson.M{
        "achievement_references_id": achievementReferencesID,
        "file_name":                 fileName,
        "folder":                    folder,
        "created_at":                time.Now(),
    }

    _, err := collection.InsertOne(context.Background(), doc)
    if err != nil {
        return "", err
    }

    return folder, nil
}

func GetAttachmentsByReferenceID(
	ctx context.Context,
	referenceID string,
) ([]model.Attachment, error) {

	collection := database.MongoDb.Collection("achievement_attachments")

	cursor, err := collection.Find(
		ctx,
		bson.M{"achievement_references_id": referenceID},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var attachments []model.Attachment

	for cursor.Next(ctx) {
		var raw bson.M
		if err := cursor.Decode(&raw); err != nil {
			return nil, err
		}

		fileName, _ := raw["file_name"].(string)

		var uploadedAt time.Time
		if dt, ok := raw["created_at"].(primitive.DateTime); ok {
			uploadedAt = dt.Time()
		}

		attachments = append(attachments, model.Attachment{
			FileName: fileName,
			FileURL: fmt.Sprintf(
				"/uploads/achievements/%s/%s",
				referenceID,
				fileName,
			),
			FileType:   filepath.Ext(fileName),
			UploadedAt: uploadedAt,
		})
	}

	return attachments, nil
}


func FindAchievementByID(ctx context.Context, id string) (model.Achievement, error) {
	collection := database.MongoDb.Collection("achievements")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Achievement{}, err
	}

	var achievement model.Achievement
	err = collection.FindOne(
		ctx,
		bson.M{"_id": oid},
	).Decode(&achievement)

	if err != nil {
		return model.Achievement{}, err
	}

	return achievement, nil
}


func (r *AchievementRepo) Delete(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}


func UpdateAchievementFieldsByID(
	ctx context.Context,
	mongoID string,
	fields bson.M,
) error {

	collection := database.MongoDb.Collection("achievements")

	oid, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return err
	}

	_, err = collection.UpdateByID(
		ctx,
		oid,
		bson.M{"$set": fields},
	)
	return err
}
