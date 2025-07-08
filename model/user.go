package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
	Email    string             `json:"email" bson:"email"`
	Fullname string             `json:"fullname" bson:"fullname"`
	IsVerified bool             `json:"is_verified" bson:"is_verified"`
}

