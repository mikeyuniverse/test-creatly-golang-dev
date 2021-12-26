package repo

import (
	"context"
	"creatly-task/internal/mongodb"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
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

func (t *TokenStorage) GetUserIDByToken(token string) (int64, error) {
	result := t.db.FindOne(context.TODO(), bson.M{"token": token})
	fmt.Println(result)
	return 0, nil
}
