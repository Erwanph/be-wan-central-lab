package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty"`
	Name             string              `bson:"name"`
	Email            string              `bson:"email"`
	Password         string              `bson:"password"`
	CreatedAt        time.Time           `bson:"created_at"`
	UpdatedAt        *time.Time          `bson:"updated_at"`
	SecretKey        string              `bson:"secret_key"`
	OTP              string              `bson:"otp"`
	OTPExpiresAt     time.Time           `bson:"otp_expires_at"`
	IsEmailVerified  bool                `bson:"is_email_verified"`
}
