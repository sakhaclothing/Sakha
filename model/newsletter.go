package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewsletterSubscription represents a user's newsletter subscription
type NewsletterSubscription struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	IsActive  bool               `json:"is_active" bson:"is_active"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// NewsletterEmail represents an email sent to subscribers
type NewsletterEmail struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Subject   string             `json:"subject" bson:"subject"`
	Content   string             `json:"content" bson:"content"`
	SentAt    time.Time          `json:"sent_at" bson:"sent_at"`
	SentCount int                `json:"sent_count" bson:"sent_count"`
	Status    string             `json:"status" bson:"status"` // "pending", "sent", "failed"
}

// NewProductNotification represents a notification for new products
type NewProductNotification struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Product   Product            `json:"product" bson:"product"`
	SentAt    time.Time          `json:"sent_at" bson:"sent_at"`
	SentCount int                `json:"sent_count" bson:"sent_count"`
	Status    string             `json:"status" bson:"status"` // "pending", "sent", "failed"
}
