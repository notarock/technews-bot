package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollections struct {
	Guild *mongo.Collection
}

var Client *mongo.Client
var collections MongoCollections

type MongodbConfig struct {
	Uri    string
	DbName string
}

func Connect(mongodbConfig MongodbConfig) (err error) {
	clientOptions := options.Client().ApplyURI(mongodbConfig.Uri)
	Client, err = mongo.NewClient(clientOptions)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = Client.Connect(ctx)
	if err != nil {
		return err
	}

	db := Client.Database(mongodbConfig.DbName)
	guildCollection := db.Collection(GUILD_COLLECTION)
	collections = MongoCollections{
		Guild: guildCollection,
	}

	return nil
}
