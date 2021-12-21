package repo

import (
	"creatly-task/internal/mongodb"

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
