package application

import (
	"context"
	"fmt"
	connectionService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/connection"
	messageService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/message"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post-microservice/model"
	"post-microservice/startup/config"
	"time"
)

type CommentService struct {
	store            model.CommentStore
	postStore        model.PostStore
	connectionClient connectionService.ConnectionServiceClient
	messageClient    messageService.MessageServiceClient
}

func NewCommentService(store model.CommentStore, postStore model.PostStore, c *config.Config) *CommentService {
	return &CommentService{
		store:            store,
		postStore:        postStore,
		messageClient:    services.NewMessageClient(fmt.Sprintf("%s:%s", c.MessageServiceHost, c.MessageServicePort)),
		connectionClient: services.NewConnectionClient(fmt.Sprintf("%s:%s", c.ConnectionServiceHost, c.ConnectionServicePort)),
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
	retVal, err := service.store.Create(ctx, comment)

	if err != nil {
		return nil, err
	}

	service.sendNotification(ctx, comment)

	return retVal, nil
}

func (service *CommentService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return service.store.Delete(ctx, id)
}

func (service *CommentService) sendNotification(ctx context.Context, comment *model.Comment) error {

	id, err := primitive.ObjectIDFromHex(comment.PostId)
	if err != nil {
		return err
	}
	post, err := service.postStore.Get(ctx, id)
	if err != nil {
		return err
	}

	connection, err := service.connectionClient.GetConnection(ctx, &connectionService.Connection{UserId: post.UserId, ConnectedUserId: comment.UserId})

	if err == nil && connection.IsCommentNotificationEnabled {
		service.messageClient.CreateNotification(ctx, &messageService.NewNotificationRequest{NotificationType: 3, Notification: &messageService.Notification{UserId: post.UserId, FromUserId: comment.UserId}})
	}
	return nil
}
