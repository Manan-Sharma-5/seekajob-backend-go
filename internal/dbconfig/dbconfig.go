package db

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
    clientInstance     *mongo.Client
    clientInstanceErr  error
    mongoOnce          sync.Once
    mongoURI           = "mongodb://localhost:27017" // Replace with env/config
    connectionTimeout  = 10 * time.Second
)

// GetMongoClient returns a singleton MongoDB client
func GetMongoClient() (*mongo.Database, error) {
    mongoOnce.Do(func() {
        ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
        defer cancel()

        clientOptions := options.Client().ApplyURI(mongoURI)
        client, err := mongo.Connect(ctx, clientOptions)
        if err != nil {
            clientInstanceErr = err
            return
        }

        // Ping to verify connection
        if err := client.Ping(ctx, nil); err != nil {
            clientInstanceErr = err
            return
        }

        clientInstance = client
        log.Println("✅ Connected to MongoDB")
    })

    return clientInstance.Database("seek-a-job"), clientInstanceErr
}
