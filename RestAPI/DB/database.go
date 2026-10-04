package database

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ConnectDB loads the MongoDB URL, verifies the connection, and returns the
// collection used by the task API. The client is returned so main can close it
// when the process shuts down.
func ConnectDB() (*mongo.Client, *mongo.Collection, error) {
	// Loading .env is convenient locally; deployed environments can provide the
	// same setting directly without needing a .env file.
	_ = godotenv.Load()

	mongoDBURL := os.Getenv("MONGODB_URL")

	if mongoDBURL == "" {
		return nil, nil, os.ErrNotExist
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	client, err := mongo.Connect(
		options.Client().ApplyURI(mongoDBURL),
	)

	if err != nil {
		return nil, nil, err
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		_ = client.Disconnect(context.Background())
		return nil, nil, err
	}

	return client, client.Database("task_api").Collection("tasks"), nil
}
