package connection

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectMongoDB() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mongoURL := "mongodb://localhost:27018"
	dbName := "Telegram"

	if v := os.Getenv("APP_DB_MONGO_URL"); v != "" {
		mongoURL = v
	}
	if v := os.Getenv("APP_DB_MONGO_DATABASE"); v != "" {
		dbName = v
	}

	// NEW KEY Keep it simple
	if v := os.Getenv("MONGO_URL"); v != "" {
		mongoURL = v
	}
	if v := os.Getenv("MONGO_DB"); v != "" {
		dbName = v
	}

	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(mongoURL).
		SetServerSelectionTimeout(10*time.Second).
		SetConnectTimeout(10*time.Second))
	if err != nil {
		fmt.Println("Failed to Initialize MongoDB")
		return nil, err
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Nearest()); err != nil {
		fmt.Println("Failed to connect MongoDB")
		return nil, err
	}

	return client.Database(dbName), nil
}
