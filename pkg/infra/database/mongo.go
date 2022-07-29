package database

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CONNECTION_STRING = "mongodb://localhost:27017/"
	DB_NAME           = "friday"
)

const MONGO_QUERY_TIMEOUT = 20 * time.Second

type MongoRepository struct {
	collection *mongo.Collection
}

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mongoInstance *MongoInstance
var mongoInstanceError error
var mongoOnce sync.Once

func New() (mongoInstance *MongoInstance, err error) {
	const mongoURI = CONNECTION_STRING + DB_NAME
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(DB_NAME)

	if err != nil {
		return
	}

	mongoInstance = &MongoInstance{Client: client, Db: db}

	return
}

func GetMongoInstance() (*MongoInstance, error) {
	mongoOnce.Do(func() {
		mongoInstance, mongoInstanceError = New()
	})

	return mongoInstance, mongoInstanceError
}
