package database

import (
    "context"
    "fmt"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
    "github.com/sirupsen/logrus"
)

// NewMongoConnection establishes a connection to MongoDB and returns the database handle.
func NewMongoConnection(logger *logrus.Logger) (*mongo.Database, error) {
    // Get connection string from environment (fallback to local default)
    uri := os.Getenv("MONGODB_URI")
    if uri == "" {
        uri = "mongodb://localhost:27017"
    }

    dbName := os.Getenv("DB_NAME")
    if dbName == "" {
        dbName = "go_api_db"
    }

    clientOps := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(clientOps)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }

    // Verify connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := client.Ping(ctx, nil); err != nil {
        return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
    }

    logger.Info("Connected to MongoDB")
    return client.Database(dbName), nil
}