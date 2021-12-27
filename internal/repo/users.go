package repo

import (
	"context"
	"creatly-task/internal/models"
	"creatly-task/internal/mongodb"
	"errors"
	"fmt"

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

	_, err = u.db.InsertOne(context.TODO(), input)
	if err != nil {
		return err
	}

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

func (u *UserStorage) GetUserByCreds(email string) (*models.UserSignInOutput, error) {
	result := u.db.FindOne(context.TODO(), bson.M{"email": email})

	if result.Err() != nil {
		return nil, result.Err()
	}

	if result.Err() == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}

	var user models.UserSignInOutput
	err := result.Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("decode error: %s", err.Error())
	}

	return &user, nil
}
