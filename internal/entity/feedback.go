package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feedback struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Message   string             `bson:"message"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt *time.Time         `bson:"updated_at"`
}
