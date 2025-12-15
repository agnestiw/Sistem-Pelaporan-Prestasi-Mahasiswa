package mongo

import (
	"context"
	"sistem-prestasi/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCompetitionLevelDistributionMongo(mongoIDs []string) ([]map[string]interface{}, error) {

	ctx := context.TODO()
	collection := database.MongoDb.Collection("achievements")

	var objectIDs []primitive.ObjectID
	for _, id := range mongoIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			objectIDs = append(objectIDs, objID)
		}
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": bson.M{"$in": objectIDs},
				"achievementType": "competition",
				"details.competitionLevel": bson.M{"$exists": true, "$ne": ""},
			},
		},
		{
			"$group": bson.M{
				"_id":   "$details.competitionLevel",
				"total": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []map[string]interface{}
	for cursor.Next(ctx) {
		var row bson.M
		cursor.Decode(&row)

		result = append(result, map[string]interface{}{
			"level": row["_id"],
			"total": row["total"],
		})
	}

	return result, nil
}
