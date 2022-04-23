package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId       string             `json:"userId"`
	Text         string             `json:"text"`
	Photo        string             `json:"photo"`
	Links        []string           `json:"links"`
	CreationDate time.Time          `json:"creationDate"`
}
