package twitter

import (
	"context"

	"github.com/mendesbarreto/friday/pkg/infra/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const TWEETS_COLLECTION = "tweets"

type TweetsRepository interface {
	Count() (int64, error)
	FindAll() ([]*Tweet, error)
	Creat(tweet *Tweet) error
	FindById(id primitive.ObjectID) (*Tweet, error)
}

type MongoTweetCollection struct {
	collection *mongo.Collection
}

func (u *MongoTweetCollection) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), database.MONGO_QUERY_TIMEOUT)

	defer cancel()

	return u.collection.CountDocuments(ctx, bson.M{})
}

func (u *MongoTweetCollection) FindAll() (*[]Tweet, error) {
	var result []Tweet

	ctx, cancel := context.WithTimeout(context.Background(), database.MONGO_QUERY_TIMEOUT)

	defer cancel()

	cursor, err := u.collection.Find(ctx, options.Find())

	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		tweet := Tweet{}

		if err := cursor.Decode(tweet); err != nil {
			return nil, err
		}

		result = append(result, tweet)
	}

	return &result, nil
}
