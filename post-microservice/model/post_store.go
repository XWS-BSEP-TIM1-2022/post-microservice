package model

import "context"

type PostStore interface {
	GetAll(ctx context.Context) ([]*Post, error)
}
