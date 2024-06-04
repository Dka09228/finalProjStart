// entity/post.go
package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	CreatedAt time.Time          `json:"created_at"`
}
