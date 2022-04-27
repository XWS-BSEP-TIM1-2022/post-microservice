package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*Comment, error)
	GetAll(ctx context.Context) ([]*Comment, error)
	GetAllFromPost(ctx context.Context, postId string) ([]*Comment, error)
	Create(ctx context.Context, comment *Comment) (*Comment, error)
	//Update(ctx context.Context, id primitive.ObjectID, comment *Comment) (*Comment, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}
