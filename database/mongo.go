// database/mongo.go
package database

import (
    "context"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() *mongo.Client {
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
    if err != nil {
        log.Fatal("MongoDB gagal connect:", err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("MongoDB tidak respons:", err)
    }

    log.Println("MongoDB berhasil terkoneksi")
    return client
}