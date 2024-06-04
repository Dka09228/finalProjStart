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
	MongoURI       string = "mongodb+srv://alisheribraev03:m3Zi0rAnHCMJpm0f@alish.ahqrfiy.mongodb.net/?retryWrites=true&w=majority&appName=Alish"
	DatabaseName   string = "golang-test"
	CollectionName string = "premleague-test"
	CounterName    string = "counters" //test
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
