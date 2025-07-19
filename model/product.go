package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	Category    string             `json:"category" bson:"category"`
	Stock       int                `json:"stock" bson:"stock"`
	IsFeatured  bool               `json:"is_featured" bson:"is_featured"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
} 