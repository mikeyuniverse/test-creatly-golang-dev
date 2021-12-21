package mongodb

import (
	"context"
	"creatly-task/internal/config"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	DB *mongo.Database
}

func New(config *config.Repo) (*Mongo, error) {
	ctx := context.Background()

	opts := options.Client()
	opts.ApplyURI(fmt.Sprintf("mongodb://%s:%s", config.Host, config.Port))

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	db := client.Database(config.DatabaseName)

	return &Mongo{
		DB: db,
	}, nil
}
