package application

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post-microservice/model"
	"time"
)

type CommentService struct {
	store model.CommentStore
}

func NewCommentService(store model.CommentStore) *CommentService {
	return &CommentService{
		store: store,
	}
}

func (service *CommentService) Get(ctx context.Context, id primitive.ObjectID) (*model.Comment, error) {
	return service.store.Get(ctx, id)
}

func (service *CommentService) GetAll(ctx context.Context) ([]*model.Comment, error) {
	return service.store.GetAll(ctx)
}

func (service *CommentService) GetAllFromPost(ctx context.Context, postId string) ([]*model.Comment, error) {
	return service.store.GetAllFromPost(ctx, postId)
}

func (service *CommentService) Create(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	comment.CreationDate = time.Now()
	return service.store.Create(ctx, comment)
}

func (service *CommentService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return service.store.Delete(ctx, id)
}
