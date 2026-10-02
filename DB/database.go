package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDb() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("error loading environment: ", err)
	}
	mongoDbUrl := os.Getenv("MONGODB_URL")

	if mongoDbUrl == "" {
		log.Fatal("MONGODB_URL is not set")
	}
	log.Println("Connecting to MongoDB...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	client, err := mongo.Connect(
		options.Client().ApplyURI(mongoDbUrl),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB connected")

}
