package repository

import (
	"context"
	"finalProjStart/db"
	"finalProjStart/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type PostRepository interface {
	Save(*entity.Post) (*entity.Post, error)
	FindAll() ([]entity.Post, error)
	FindByID(int) (*entity.Post, error)
	Update(int, *entity.Post) (*entity.Post, error)
	Delete(int) error
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

// getNextSequence generates the next auto-increment ID
func (r *repo) getNextSequence() (int, error) {
	collection := r.client.Database(db.DatabaseName).Collection(db.CounterName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result struct {
		Seq int `bson:"seq"`
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	filter := bson.M{"_id": db.CollectionName}
	update := bson.M{"$inc": bson.M{"seq": 1}}

	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.Seq, nil
}

func (r *repo) Save(post *entity.Post) (*entity.Post, error) {
	collection := r.client.Database(db.DatabaseName).Collection(db.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, err := r.getNextSequence()
	if err != nil {
		log.Fatalf("Failed to get next sequence: %v", err)
		return nil, err
	}
	post.ID = id

	_, err = collection.InsertOne(ctx, post)
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

func (r *repo) FindByID(id int) (*entity.Post, error) {
	collection := r.client.Database(db.DatabaseName).Collection(db.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var post entity.Post
	filter := bson.M{"id": id}
	err := collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Fatalf("Failed to find post by ID: %v", err)
		return nil, err
	}

	return &post, nil
}

func (r *repo) Update(id int, updatedPost *entity.Post) (*entity.Post, error) {
	collection := r.client.Database(db.DatabaseName).Collection(db.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": updatedPost}

	result := collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Fatalf("Failed to update post: %v", result.Err())
		return nil, result.Err()
	}

	var post entity.Post
	err := result.Decode(&post)
	if err != nil {
		log.Fatalf("Failed to decode updated post: %v", err)
		return nil, err
	}

	return &post, nil
}

func (r *repo) Delete(id int) error {
	collection := r.client.Database(db.DatabaseName).Collection(db.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatalf("Failed to delete post: %v", err)
		return err
	}

	return nil
}
