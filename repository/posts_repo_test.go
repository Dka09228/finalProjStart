package repository

import (
	"context"
	"finalProjStart/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestRepository_Save(t *testing.T) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://alisheribraev03:m3Zi0rAnHCMJpm0f@alish.ahqrfiy.mongodb.net/?retryWrites=true&w=majority&appName=Alish"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("test_db")
	repo := NewMongoDBPostRepository(db)

	post := &entity.Post{
		Title:     "Test Post",
		Content:   "This is a test post.",
		CreatedAt: time.Now(),
	}

	_, err = repo.Save(post)
	if err != nil {
		t.Errorf("Error saving post: %v", err)
	}

}

func TestRepository_FindAll(t *testing.T) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://alisheribraev03:m3Zi0rAnHCMJpm0f@alish.ahqrfiy.mongodb.net/?retryWrites=true&w=majority&appName=Alish"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("test_db")
	repo := NewMongoDBPostRepository(db)

	_, err = repo.FindAll()
	if err != nil {
		t.Errorf("Error finding all posts: %v", err)
	}
}

func TestRepository_Update(t *testing.T) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://alisheribraev03:m3Zi0rAnHCMJpm0f@alish.ahqrfiy.mongodb.net/?retryWrites=true&w=majority&appName=Alish"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("test_db")
	repo := NewMongoDBPostRepository(db)

	post := &entity.Post{
		Title:     "Test Post",
		Content:   "This is a test post.",
		CreatedAt: time.Now(),
	}
	_, err = repo.Save(post)
	if err != nil {
		t.Fatalf("Error saving post: %v", err)
	}

	post.Content = "Updated description"
	err = repo.Update(post)
	if err != nil {
		t.Errorf("Error updating post: %v", err)
	}
}

func TestRepository_Delete(t *testing.T) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://alisheribraev03:m3Zi0rAnHCMJpm0f@alish.ahqrfiy.mongodb.net/?retryWrites=true&w=majority&appName=Alish"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("test_db")
	repo := NewMongoDBPostRepository(db)

	post := &entity.Post{
		Title:     "Test Post",
		Content:   "This is a test post.",
		CreatedAt: time.Now(),
	}
	_, err = repo.Save(post)
	if err != nil {
		t.Fatalf("Error saving post: %v", err)
	}

	err = repo.Delete(post.ID)
	if err != nil {
		t.Errorf("Error deleting post: %v", err)
	}

	foundPost, err := repo.FindByID(post.ID)
	if err != nil {
		t.Errorf("Error finding post by ID: %v", err)
	}
	if foundPost != nil {
		t.Errorf("Post still found after deletion")
	}
}

func TestNonExistentID(t *testing.T) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://alisheribraev03:m3Zi0rAnHCMJpm0f@alish.ahqrfiy.mongodb.net/?retryWrites=true&w=majority&appName=Alish"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("test_db")
	repo := NewMongoDBPostRepository(db)

	nonExistentID := primitive.NewObjectID()

	foundPost, err := repo.FindByID(nonExistentID)
	if err != nil {
		t.Errorf("Error finding post by non-existent ID: %v", err)
	}
	if foundPost != nil {
		t.Errorf("Expected nil post for non-existent ID, got non-nil")
	}
}
