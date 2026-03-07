package database

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollections struct {
	Guild   *mongo.Collection
	Article *mongo.Collection
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = Client.Connect(ctx)
	if err != nil {
		return err
	}

	db := Client.Database(mongodbConfig.DbName)

	guildCollection := db.Collection(GUILD_COLLECTION)
	articleCollection := db.Collection(ARTICLE_COLLECTION)

	collections = MongoCollections{
		Guild:   guildCollection,
		Article: articleCollection,
	}

	return nil
}

func Healthcheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collections.Guild.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	logrus.Infof("Serving %d guilds", count)
	return nil
}
