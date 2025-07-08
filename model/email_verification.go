package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmailVerification struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	OTP       string             `json:"otp" bson:"otp"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	Used      bool               `json:"used" bson:"used"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
} 