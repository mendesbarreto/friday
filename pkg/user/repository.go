package user

import (
	"context"

	"github.com/mendesbarreto/friday/pkg/infra/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const USER_COLLECTION = "users"

type UserRepository interface {
	FindAll() ([]*User, error)
	Create(*User) (*User, error)
}

type MongoUserRepository struct {
	collection *mongo.Collection
}

func (u *MongoUserRepository) FindByUserName(username string) (*User, error) {
	var result *User

	ctx, cancel := context.WithTimeout(context.Background(), database.MONGO_QUERY_TIMEOUT)

	defer cancel()

	err := u.collection.FindOne(ctx, bson.M{"username": username}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *MongoUserRepository) FindAll() ([]*User, error) {
	var result []*User

	ctx, cancel := context.WithTimeout(context.Background(), database.MONGO_QUERY_TIMEOUT)

	defer cancel()

	cursor, err := u.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "username", Value: 1}}))

	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		user := new(User)

		if err := cursor.Decode(user); err != nil {
			return nil, err
		}

		result = append(result, user)
	}

	return result, cursor.Err()
}

func (u *MongoUserRepository) Create(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), database.MONGO_QUERY_TIMEOUT)

	defer cancel()

	_, err := u.collection.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (u *MongoUserRepository) FindById() (*User, error) {
	return nil, nil
}

func NewUserRepository() (*MongoUserRepository, error) {
	mongoInstance, err := database.GetMongoInstance()

	if err != nil {
		return nil, err
	}

	return &MongoUserRepository{collection: mongoInstance.Db.Collection(USER_COLLECTION)}, nil
}
