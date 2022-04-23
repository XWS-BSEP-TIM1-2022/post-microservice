package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*Post, error)
	GetAll(ctx context.Context) ([]*Post, error)
	GetAllFromUser(ctx context.Context, userId string) ([]*Post, error)
	Create(ctx context.Context, post *Post) (*Post, error)
	Update(ctx context.Context, id primitive.ObjectID, post *Post) (*Post, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}
