package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password,omitempty" bson:"password,omitempty"`
	Role         string             `json:"role" bson:"role"`
	RegisteredAt time.Time          `json:"registered_at" bson:"registered_at"`
}
