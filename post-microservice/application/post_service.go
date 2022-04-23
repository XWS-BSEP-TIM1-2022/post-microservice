package application

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post-microservice/model"
	"time"
)

type PostService struct {
	store model.PostStore
}

func NewPostService(store model.PostStore) *PostService {
	return &PostService{
		store: store,
	}
}

func (service *PostService) Get(ctx context.Context, id primitive.ObjectID) (*model.Post, error) {
	return service.store.Get(ctx, id)
}

func (service *PostService) GetAll(ctx context.Context) ([]*model.Post, error) {
	return service.store.GetAll(ctx)
}

func (service *PostService) GetAllFromUser(ctx context.Context, userId string) ([]*model.Post, error) {
	return service.store.GetAllFromUser(ctx, userId)
}

func (service *PostService) Create(ctx context.Context, post *model.Post) (*model.Post, error) {
	post.CreationDate = time.Now()

	return service.store.Create(ctx, post)
}

func (service *PostService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return service.store.Delete(ctx, id)
}
