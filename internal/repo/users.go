package repo

import (
	"context"
	"creatly-task/internal/models"
	"creatly-task/internal/mongodb"

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
	u.db.InsertOne(context.TODO(), input)
	return nil
}

func (u *UserStorage) GetUserByCreds(input *models.UserSignInInput) error {
	u.db.InsertOne(context.TODO(), input)
	return nil
}
