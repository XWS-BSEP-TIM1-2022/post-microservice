package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Comment struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId       string             `json:"userId" bson:"userId"`
	PostId       string             `json:"postId" bson:"postId"`
	Text         string             `json:"text"`
	CreationDate time.Time          `json:"creationDate"`
}
