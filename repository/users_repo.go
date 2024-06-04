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

type UserRepository interface {
	SaveUser(user *entity.User) error
	FindUserByEmail(email string) (*entity.User, error)
	FindAllUsers() ([]entity.User, error)
}

type MongoDBUserRepository struct {
	collection *mongo.Collection
}

func NewMongoDBUserRepository(database *mongo.Database) UserRepository {
	collection := database.Collection("users")
	return &MongoDBUserRepository{
		collection: collection,
	}
}

func (repo *MongoDBUserRepository) SaveUser(user *entity.User) error {
	user.RegisteredAt = time.Now()

	insertResult, err := repo.collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("Error inserting user: %v\n", err)
		return err
	}

	if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}

	return nil
}

func (repo *MongoDBUserRepository) FindUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	filter := bson.D{{"email", email}}

	err := repo.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		log.Printf("Error finding user by email: %v\n", err)
		return nil, err
	}

	return &user, nil
}

func (repo *MongoDBUserRepository) FindAllUsers() ([]entity.User, error) {
	var users []entity.User

	cur, err := repo.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Printf("Error finding all users: %v\n", err)
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var user entity.User
		err := cur.Decode(&user)
		if err != nil {
			log.Printf("Error decoding user: %v\n", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		log.Printf("Cursor error: %v\n", err)
		return nil, err
	}

	return users, nil
}
