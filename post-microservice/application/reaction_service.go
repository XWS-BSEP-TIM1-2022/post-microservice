package application

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post-microservice/model"
)

type ReactionService struct {
	store model.ReactionStore
}

func NewReactionService(store model.ReactionStore) *ReactionService {
	return &ReactionService{
		store: store,
	}
}

func (service *ReactionService) Get(ctx context.Context, id primitive.ObjectID) (*model.Reaction, error) {
	Log.Info("Get reaction by id: " + id.Hex())
	return service.store.Get(ctx, id)
}

func (service *ReactionService) GetAll(ctx context.Context) ([]*model.Reaction, error) {
	Log.Info("Get all reactions")
	return service.store.GetAll(ctx)
}

func (service *ReactionService) GetAllFromPost(ctx context.Context, postId string) ([]*model.Reaction, error) {
	Log.Info("Get all reactions of post with id: " + postId)
	return service.store.GetAllFromPost(ctx, postId)
}

func (service *ReactionService) Create(ctx context.Context, reaction *model.Reaction) (*model.Reaction, error) {
	Log.Info("Create reaction")
	return service.store.Create(ctx, reaction)
}

func (service *ReactionService) Delete(ctx context.Context, id primitive.ObjectID) error {
	Log.Info("Delete reaction by id: " + id.Hex())
	return service.store.Delete(ctx, id)
}
