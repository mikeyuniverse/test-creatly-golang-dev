package repo

import (
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
