package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    Id primitive.ObjectID `json:"id" bson:"_id,omitempty"` 
    Username string `json:"title" bson:"title"`
    Password string `json:"password" bson:"password"`
    CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt" bson:"updated"`
}

