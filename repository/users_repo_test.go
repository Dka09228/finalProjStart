package repository

import (
	"context"
	"finalProjStart/entity"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestSaveUser(t *testing.T) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://alisheribraev03:m3Zi0rAnHCMJpm0f@alish.ahqrfiy.mongodb.net/?retryWrites=true&w=majority&appName=Alish"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("test_db")
	repo := NewMongoDBUserRepository(db)

	user := &entity.User{
		Email:        "test@example.com",
		Password:     "password",
		RegisteredAt: time.Now(),
	}

	err = repo.SaveUser(user)
	if err != nil {
		t.Errorf("Error saving user: %v", err)
	}
}

func TestFindUserByEmail(t *testing.T) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://alisheribraev03:m3Zi0rAnHCMJpm0f@alish.ahqrfiy.mongodb.net/?retryWrites=true&w=majority&appName=Alish"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("test_db")
	repo := NewMongoDBUserRepository(db)

	existingUserEmail := "test@example.com"
	foundUser, err := repo.FindUserByEmail(existingUserEmail)
	if err != nil {
		t.Errorf("Error finding user by email: %v", err)
	}
	if foundUser == nil {
		t.Errorf("Expected user not found")
	}
}

func TestFindAllUsers(t *testing.T) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://alisheribraev03:m3Zi0rAnHCMJpm0f@alish.ahqrfiy.mongodb.net/?retryWrites=true&w=majority&appName=Alish"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("test_db")
	repo := NewMongoDBUserRepository(db)

	_, err = repo.FindAllUsers()
	if err != nil {
		t.Errorf("Error finding all users: %v", err)
	}
}

func TestNonExistentEmail(t *testing.T) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://alisheribraev03:m3Zi0rAnHCMJpm0f@alish.ahqrfiy.mongodb.net/?retryWrites=true&w=majority&appName=Alish"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("test_db")
	repo := NewMongoDBUserRepository(db)

	nonExistentEmail := "nonexistent@example.com"
	foundUser, err := repo.FindUserByEmail(nonExistentEmail)
	if err == nil {
		t.Error("Expected error due to non-existent email, got nil")
	}
	if foundUser != nil {
		t.Error("Expected nil user for non-existent email, got non-nil")
	}
}
