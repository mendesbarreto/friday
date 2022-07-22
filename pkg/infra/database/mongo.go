package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

const dbName = "friday"
const mongoURI = "mongodb://localhost:27017/" + dbName

func New() (mongoInstance MongoInstance, err error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		return
	}

	mongoInstance = MongoInstance{Client: client, Db: db}

	return
}
