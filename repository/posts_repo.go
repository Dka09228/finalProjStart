package repository

import (
	"context"
	"finalProjStart/db"
	"finalProjStart/entity"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository interface {
	Save(*entity.Post) (*entity.Post, error)
	FindAll() ([]entity.Post, error)
}

type repo struct {
	client *mongo.Client
}

// NewPostRepository creates a new PostRepository with MongoDB client
func NewPostRepository() (PostRepository, error) {
	client, err := db.ConnectMongoDB()
	if err != nil {
		return nil, err
	}

	return &repo{client: client}, nil
}

func (r *repo) Save(post *entity.Post) (*entity.Post, error) {
	collection := r.client.Database(db.DatabaseName).Collection(db.CollectionName)
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
	collection := r.client.Database(db.DatabaseName).Collection(db.CollectionName)
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
