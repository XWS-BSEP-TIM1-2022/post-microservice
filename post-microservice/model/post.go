package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId       string             `json:"userId" bson:"userId"`
	Text         string             `json:"text"`
	Image        string             `json:"image"`
	Links        []string           `json:"links"`
	CreationDate time.Time          `json:"creationDate"`
}
