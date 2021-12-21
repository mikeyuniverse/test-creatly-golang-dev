package mongodb

import (
	"context"
	"creatly-task/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	DB *mongo.Database
}

func New(config *config.Repo) (*Mongo, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	db := client.Database(config.DatabaseName)
	return &Mongo{
		DB: db,
	}, nil
}
