// repository/post_repository.go
package repository

import (
	"context"
	"errors"
	"finalProjStart/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type PostRepository interface {
	Save(post *entity.Post) (*entity.Post, error)
	FindAll() ([]entity.Post, error)
	FindByID(id primitive.ObjectID) (*entity.Post, error)
	Update(post *entity.Post) error
	Delete(id primitive.ObjectID) error
}

type MongoDBPostRepository struct {
	collection *mongo.Collection
}

func NewMongoDBPostRepository(database *mongo.Database) PostRepository {
	collection := database.Collection("posts")
	return &MongoDBPostRepository{
		collection: collection,
	}
}

func (repo *MongoDBPostRepository) Save(post *entity.Post) (*entity.Post, error) {
	post.CreatedAt = time.Now()

	_, err := repo.collection.InsertOne(context.TODO(), post)
	if err != nil {
		log.Printf("Error inserting post: %v\n", err)
		return nil, err
	}

	return post, nil
}

func (repo *MongoDBPostRepository) FindAll() ([]entity.Post, error) {
	var posts []entity.Post

	cur, err := repo.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Printf("Error finding all posts: %v\n", err)
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var post entity.Post
		err := cur.Decode(&post)
		if err != nil {
			log.Printf("Error decoding post: %v\n", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := cur.Err(); err != nil {
		log.Printf("Cursor error: %v\n", err)
		return nil, err
	}

	return posts, nil
}

func (repo *MongoDBPostRepository) FindByID(id primitive.ObjectID) (*entity.Post, error) {
	var post entity.Post
	filter := bson.D{{"_id", id}}

	err := repo.collection.FindOne(context.TODO(), filter).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		log.Printf("Error finding post by ID: %v\n", err)
		return nil, err
	}

	return &post, nil
}

func (repo *MongoDBPostRepository) Update(post *entity.Post) error {
	filter := bson.D{{"_id", post.ID}}
	update := bson.D{{"$set", post}}

	_, err := repo.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Error updating post: %v\n", err)
		return err
	}

	return nil
}

func (repo *MongoDBPostRepository) Delete(id primitive.ObjectID) error {
	filter := bson.D{{"_id", id}}

	_, err := repo.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Error deleting post: %v\n", err)
		return err
	}

	return nil
}
