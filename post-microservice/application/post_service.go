package application

import (
	"context"
	"post-microservice/model"
)

type PostService struct {
	store model.PostStore
}

func NewPostService(store model.PostStore) *PostService {
	return &PostService{
		store: store,
	}
}

func (service *PostService) GetAll(ctx context.Context) ([]*model.Post, error) {
	return service.store.GetAll(ctx)
}
