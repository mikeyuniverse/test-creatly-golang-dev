package repo

import (
	"context"
	"creatly-task/internal/models"
	"creatly-task/internal/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FilesRepo struct {
	db *mongo.Collection
}

func newFilesRepo(mongo *mongodb.Mongo, collectionName string) *FilesRepo {
	collection := mongo.DB.Collection(collectionName)
	return &FilesRepo{
		db: collection,
	}
}

func (f *FilesRepo) All() ([]models.FileOut, error) {
	cursor, err := f.db.Find(context.TODO(), bson.M{})
	if err != nil {
		return []models.FileOut{}, err
	}

	results := make([]models.FileOut, 1)
	err = cursor.All(context.Background(), &results)
	if err != nil {
		return []models.FileOut{}, err
	}

	return results, nil

}
