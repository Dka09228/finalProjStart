package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DatabaseName string = "golangInterview"
	DatabaseName   string = "golang-test"
	CollectionName string = "premleague-test"
	CounterName    string = "counters"
)

func ConnectMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(MongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	// Check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return client, nil
}
