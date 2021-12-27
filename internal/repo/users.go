package repo

import (
	"context"
	"creatly-task/internal/models"
	"creatly-task/internal/mongodb"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStorage struct {
	db *mongo.Collection
}

func newUsersRepo(mongo *mongodb.Mongo, collectionName string) *UserStorage {
	colection := mongo.DB.Collection(collectionName)
	return &UserStorage{
		db: colection,
	}
}

func (u *UserStorage) CreateUser(input *models.UserSignUpInput) error {
	err := u.checkUserExists(input.Email)
	if err != nil {
		return err
	}

	u.db.InsertOne(context.TODO(), input)
	return nil
}

func (u *UserStorage) checkUserExists(email string) error {
	result := u.db.FindOne(context.TODO(), bson.M{"email": email})

	if result.Err() == mongo.ErrNoDocuments {
		// User with this email not exists
		return nil
	}

	var dbUser models.UserSignInInput
	err := result.Decode(&dbUser)
	if err != nil {
		return err
	}

	if dbUser.Email != "" {
		// User with this email exists
		return errors.New("user exists")
	}

	// TODO Unknown logic: whats doing?
	return nil
}
