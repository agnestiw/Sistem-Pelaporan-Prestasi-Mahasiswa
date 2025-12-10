// app/repository/mongo/achievement_mongo.go
package mongo

import (
	"context"
	model "sistem-prestasi/app/model/mongo"
	"sistem-prestasi/database"

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

	return result.InsertedID.(interface{}).(string), nil
}

func DeleteAchievement(ctx context.Context, mongoID string) error {
	collection := database.MongoDb.Collection("achievements")

	_, err := collection.DeleteOne(ctx, bson.M{
		"_id": mongoID,
	})

	return err
}

func FindAchievementByID(ctx context.Context, collection *mongo.Collection, id string) (*model.Achievement, error) {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var achievement model.Achievement
	err = collection.FindOne(
		ctx,
		bson.M{"_id": oid},
	).Decode(&achievement)

	if err != nil {
		return nil, err
	}

	return &achievement, nil
}

func (r *AchievementRepo) Delete(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func UpdateAchievementByID(
	ctx context.Context,
	mongoID string,
	input model.Achievement,
) error {

	collection := database.MongoDb.Collection("achievements")

	objectID, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": input,
	}

	_, err = collection.UpdateByID(ctx, objectID, update)
	return err
}
