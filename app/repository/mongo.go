package repository

import "go.mongodb.org/mongo-driver/mongo"

var MongoClient *mongo.Client

func InitMongo(client *mongo.Client) {
    MongoClient = client
}
