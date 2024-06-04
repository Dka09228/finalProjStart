package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestScrapeAndStore(t *testing.T) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://alisheribraev03:m3Zi0rAnHCMJpm0f@alish.ahqrfiy.mongodb.net/?retryWrites=true&w=majority&appName=Alish"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	repo := NewScrapeRepository(client, nil)

	url := "http://example.com"
	tableID := "statsTable"
	collectionName := "teamStats"
	err = repo.ScrapeAndStoreData(url, tableID, collectionName)
	if err != nil {
		t.Errorf("Error scraping and storing data: %v", err)
	}
}
