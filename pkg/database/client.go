package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

type MongodbConfig struct {
	Uri string
}

func Connect(mongodbConfig MongodbConfig) error {
	clientOptions := options.Client().ApplyURI(mongodbConfig.Uri)
	Client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = Client.Connect(ctx)
	if err != nil {
		return err
	}

	return nil
}