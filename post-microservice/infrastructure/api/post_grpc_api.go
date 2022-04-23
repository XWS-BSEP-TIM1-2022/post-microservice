package api

import (
	"context"
	"errors"
	postService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/post"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post-microservice/application"
)

type PostHandler struct {
	postService.UnimplementedPostServiceServer
	service *application.PostService
}

func NewPostHandler(service *application.PostService) *PostHandler {
	return &PostHandler{
		service: service,
	}
}

func (handler *PostHandler) GetRequest(ctx context.Context, in *postService.PostIdRequest) (*postService.PostResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	id := in.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	post, err := handler.service.Get(ctx, objectId)
	if err != nil {
		return nil, err
	}
	postPb := mapPost(post)
	response := &postService.PostResponse{
		Post: postPb,
	}
	return response, nil
}

func (handler *PostHandler) GetAllRequest(ctx context.Context, in *postService.EmptyRequest) (*postService.PostsResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	posts, err := handler.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	response := &postService.PostsResponse{
		Posts: []*postService.Post{},
	}
	for _, post := range posts {
		current := mapPost(post)
		response.Posts = append(response.Posts, current)
	}
	return response, nil
}

func (handler *PostHandler) GetAllFromUserRequest(ctx context.Context, in *postService.UserPostsRequest) (*postService.PostsResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	userId := in.UserId
	posts, err := handler.service.GetAllFromUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	response := &postService.PostsResponse{
		Posts: []*postService.Post{},
	}
	for _, post := range posts {
		current := mapPost(post)
		response.Posts = append(response.Posts, current)
	}
	return response, nil
}

func (handler *PostHandler) CreateRequest(ctx context.Context, in *postService.PostRequest) (*postService.PostResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "CreateRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	if in.Post.Text == "" || in.Post.UserId == "" {
		return nil, errors.New("not entered required fields")
	}

	postFromRequest := mapPostPb(in.Post)
	post, err := handler.service.Create(ctx, postFromRequest)
	if err != nil {
		return nil, err
	}
	postPb := mapPost(post)
	response := &postService.PostResponse{
		Post: postPb,
	}
	return response, nil
}

func (handler *PostHandler) DeleteRequest(ctx context.Context, in *postService.PostIdRequest) (*postService.EmptyRequest, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "DeleteRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	id, _ := primitive.ObjectIDFromHex(in.Id)
	handler.service.Delete(ctx, id)
	response := &postService.EmptyRequest{}
	return response, nil
}
