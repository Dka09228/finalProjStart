package repository

import (
	"context"
	"finalProjStart/entity"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type PostRepository interface {
	Save(*entity.Post) (*entity.Post, error)
	FindAll() ([]entity.Post, error)
}

type repo struct {
	client *mongo.Client
}

const (
	databaseName   string = "golangInterview"
	collectionName string = "posts"
	mongoURI       string = "mongodb://localhost:27017"
)

// NewPostRepository creates a new PostRepository with MongoDB client
func NewPostRepository() (PostRepository, error) {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return &repo{client: client}, nil
}

func (r *repo) Save(post *entity.Post) (*entity.Post, error) {
	collection := r.client.Database(databaseName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, post)
	if err != nil {
		log.Fatalf("Failed to insert a new post: %v", err)
		return nil, err
	}

	return post, nil
}

func (r *repo) FindAll() ([]entity.Post, error) {
	collection := r.client.Database(databaseName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to find posts: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []entity.Post
	for cursor.Next(ctx) {
		var post entity.Post
		err := cursor.Decode(&post)
		if err != nil {
			log.Fatalf("Failed to decode post: %v", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("Cursor error: %v", err)
		return nil, err
	}

	return posts, nil
}
