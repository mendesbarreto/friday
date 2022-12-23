package tweets

import (
	"context"

	"github.com/mendesbarreto/friday/pkg/infra/database"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


const TWEETS_COLLECTION = "tweets"

type Tweet struct {
    ID primitive.ObjectID `json:"id" bson:"_id,omitempty"` 
}

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
    
    return u.collection.CountDocuments(ctx, bson.M{}, bson.D{{Key: "_id", Value: 1}})
}

func (u *MongoTweetCollection) FindAll() ([]*Tweet, error) {
    var result []*Tweet

    ctx, cancel := context.WithTimeout(context.Background(), database.MONGO_QUERY_TIMEOUT)

    defer cancel()

    cursor, err := u.collection.Find(ctx, options.Find())

    if err != nil {
        reutnr nil, err
    }

    for cursor.Next(ctx) {
        tweet := new(Tweet)

        if err := cursor.Decode(tweet); err != nil {
            return nil, err
        }

        result = append(result, tweet)
    }

    return result, nil
}
