package repo

import (
	"creatly-task/internal/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

type TokenStorage struct {
	db *mongo.Collection
}

func newTokensRepo(mongo *mongodb.Mongo, collectionName string) *TokenStorage {
	collection := mongo.DB.Collection(collectionName)
	return &TokenStorage{
		db: collection,
	}
}
