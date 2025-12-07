package repository

import (
    "context"
    "os"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

func SoftDeleteMongo(client *mongo.Client, mongoID string) error {
    coll := client.Database(os.Getenv("DB_NAME")).Collection("achievements")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _, err := coll.UpdateOne(ctx, bson.M{"_id": mongoID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
    return err
}