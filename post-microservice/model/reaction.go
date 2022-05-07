package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reaction struct {
	Id     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId string             `json:"userId" bson:"userId"`
	PostId string             `json:"postId" bson:"postId"`
	Type   bool               `json:"type"`
}
