package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserSignUpInput struct {
	Email    string `json:"email",bson:"email"`
	Password string `json:"password",bson:"password"`
}

type UserSignInInput struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
}

type UserSignInOutput struct {
	UserID primitive.ObjectID `bson:"_id"`
	// UserID   string             `json:"id",bson:"_id"`
	Email    string `json:"email",bson:"email"`
	Password string `json:"password",bson:"password"`
}
