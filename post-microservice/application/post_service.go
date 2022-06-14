package application

import (
	"context"
	"fmt"
	connectionService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/connection"
	messageService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/message"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/services"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post-microservice/model"
	"post-microservice/startup/config"
	"time"
)

type PostService struct {
	store            model.PostStore
	connectionClient connectionService.ConnectionServiceClient
	messageClient    messageService.MessageServiceClient
}

var Log = logrus.New()

func NewPostService(store model.PostStore, c *config.Config) *PostService {
	return &PostService{
		store:            store,
		messageClient:    services.NewMessageClient(fmt.Sprintf("%s:%s", c.MessageServiceHost, c.MessageServicePort)),
		connectionClient: services.NewConnectionClient(fmt.Sprintf("%s:%s", c.ConnectionServiceHost, c.ConnectionServicePort)),
	}
}

func (service *PostService) Get(ctx context.Context, id primitive.ObjectID) (*model.Post, error) {
	Log.Info("Get post by id:" + id.Hex())
	return service.store.Get(ctx, id)
}

func (service *PostService) GetAll(ctx context.Context) ([]*model.Post, error) {
	Log.Info("Get all posts")
	return service.store.GetAll(ctx)
}

func (service *PostService) GetAllFromUser(ctx context.Context, userId string) ([]*model.Post, error) {
	Log.Info("Get all posts from user with id:" + userId)
	return service.store.GetAllFromUser(ctx, userId)
}

func (service *PostService) Create(ctx context.Context, post *model.Post) (*model.Post, error) {
	Log.Info("Create new Post")
	post.CreationDate = time.Now()

	retVal, err := service.store.Create(ctx, post)

	if err != nil {
		Log.Error("Error while creating post- " + err.Error())
		return nil, err
	}

	service.SendNotifications(ctx, retVal.UserId, 2)

	return retVal, nil
}

func (service *PostService) Delete(ctx context.Context, id primitive.ObjectID) error {
	Log.Info("Deleting post with id: " + id.Hex())
	return service.store.Delete(ctx, id)
}

func (service *PostService) SendNotifications(ctx context.Context, userId string, notificationType int32) error {
	Log.Info("Sending notification for post creation by user with id: " + userId)

	span := tracer.StartSpanFromContextMetadata(ctx, "SendNotifications")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	followers, err := service.connectionClient.GetFollowers(ctx, &connectionService.UserIdRequest{UserId: userId})

	if err != nil {
		Log.Error("Failed to get followers of user with id: " + userId)
		return err
	}

	for _, follower := range followers.Connections {
		if notificationType == 1 && follower.IsMessageNotificationEnabled {
			service.messageClient.CreateNotification(ctx, &messageService.NewNotificationRequest{NotificationType: notificationType, Notification: &messageService.Notification{UserId: follower.UserId, FromUserId: userId}})
		} else if notificationType == 2 && follower.IsPostNotificationEnabled {
			service.messageClient.CreateNotification(ctx, &messageService.NewNotificationRequest{NotificationType: notificationType, Notification: &messageService.Notification{UserId: follower.UserId, FromUserId: userId}})
		} else if notificationType == 3 && follower.IsCommentNotificationEnabled {
			service.messageClient.CreateNotification(ctx, &messageService.NewNotificationRequest{NotificationType: notificationType, Notification: &messageService.Notification{UserId: follower.UserId, FromUserId: userId}})
		}
	}
	return nil
}
