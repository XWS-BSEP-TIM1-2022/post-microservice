package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReactionStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*Reaction, error)
	GetAll(ctx context.Context) ([]*Reaction, error)
	GetAllFromPost(ctx context.Context, postId string) ([]*Reaction, error)
	Create(ctx context.Context, comment *Reaction) (*Reaction, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}
