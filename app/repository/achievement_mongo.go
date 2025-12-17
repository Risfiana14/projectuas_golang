package repository

import (
    "context"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "projectuas/app/model"
)

var MongoClient *mongo.Client

func InitMongo(client *mongo.Client) {
    MongoClient = client
}

func achievementsColl() *mongo.Collection {
    return MongoClient.Database(os.Getenv("DB_NAME")).Collection("achievements")
}

// CREATE
func InsertAchievementMongo(ach *model.Achievement) (string, error) {

    // ✅ SRS: attachments default empty array
    if ach.Attachments == nil {
        ach.Attachments = []model.Attachment{}
    }

    res, err := achievementsColl().InsertOne(context.TODO(), ach)
    if err != nil {
        return "", err
    }

    oid := res.InsertedID.(primitive.ObjectID).Hex()
    return oid, nil
}

// GET BY MONGO ID
func GetAchievementMongoByID(id string) (*model.Achievement, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var ach model.Achievement
    err = achievementsColl().FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&ach)
    if err != nil {
        return nil, err
    }
    return &ach, nil
}

// UPDATE
func UpdateAchievementMongo(id string, update bson.M) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    // 🔒 Guard agar attachments tidak null
    if v, ok := update["attachments"]; ok && v == nil {
        delete(update, "attachments")
    }

    _, err = achievementsColl().UpdateOne(
        context.TODO(),
        bson.M{"_id": oid},
        bson.M{"$set": update},
    )
    return err
}

// FR-005
// SOFT DELETE ACHIEVEMENT (MongoDB)
func SoftDeleteAchievementMongo(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = achievementsColl().UpdateOne(
		context.TODO(),
		bson.M{"_id": oid},
		bson.M{
			"$set": bson.M{
				"status":     "deleted",
				"deleted_at": time.Now(),
				"updated_at": time.Now(),
			},
		},
	)
	return err
}
