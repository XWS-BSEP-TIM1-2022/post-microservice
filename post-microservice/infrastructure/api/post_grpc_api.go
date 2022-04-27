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
	postService    *application.PostService
	commentService *application.CommentService
}

func NewPostHandler(postService *application.PostService, commentService *application.CommentService) *PostHandler {
	return &PostHandler{
		postService:    postService,
		commentService: commentService,
	}
}

/////////////////////////////// POSTS GRPC API ///////////////////////////////

func (handler *PostHandler) GetRequest(ctx context.Context, in *postService.PostIdRequest) (*postService.PostResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	id := in.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	post, err := handler.postService.Get(ctx, objectId)
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

	posts, err := handler.postService.GetAll(ctx)
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
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllFromUserRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	userId := in.UserId
	posts, err := handler.postService.GetAllFromUser(ctx, userId)
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
	post, err := handler.postService.Create(ctx, postFromRequest)
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
	handler.postService.Delete(ctx, id)
	response := &postService.EmptyRequest{}
	return response, nil
}

/////////////////////////////// COMMENTS GRPC API ///////////////////////////////

func (handler *PostHandler) GetCommentRequest(ctx context.Context, in *postService.CommentIdRequest) (*postService.CommentResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	id := in.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	comment, err := handler.commentService.Get(ctx, objectId)
	if err != nil {
		return nil, err
	}
	commentPb := mapComment(comment)
	response := &postService.CommentResponse{
		Comment: commentPb,
	}
	return response, nil
}

func (handler *PostHandler) GetAllCommentsRequest(ctx context.Context, in *postService.EmptyRequest) (*postService.CommentsResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	comments, err := handler.commentService.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	response := &postService.CommentsResponse{
		Comments: []*postService.Comment{},
	}
	for _, comment := range comments {
		current := mapComment(comment)
		response.Comments = append(response.Comments, current)
	}
	return response, nil
}

func (handler *PostHandler) GetAllCommentsFromPostRequest(ctx context.Context, in *postService.PostCommentsRequest) (*postService.CommentsResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllFromPostRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	postId := in.PostId
	comments, err := handler.commentService.GetAllFromPost(ctx, postId)
	if err != nil {
		return nil, err
	}
	response := &postService.CommentsResponse{
		Comments: []*postService.Comment{},
	}
	for _, comment := range comments {
		current := mapComment(comment)
		response.Comments = append(response.Comments, current)
	}
	return response, nil
}

func (handler *PostHandler) CreateCommentRequest(ctx context.Context, in *postService.CommentRequest) (*postService.CommentResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "CreateRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	postId := in.PostId
	if in.Comment.Text == "" || in.Comment.UserId == "" || postId == "" {
		return nil, errors.New("not entered required fields")
	}

	commentFromRequest := mapCommentPb(in.Comment)
	commentFromRequest.PostId = postId
	comment, err := handler.commentService.Create(ctx, commentFromRequest)
	if err != nil {
		return nil, err
	}
	commentPb := mapComment(comment)
	response := &postService.CommentResponse{
		Comment: commentPb,
	}
	return response, nil
}

func (handler *PostHandler) DeleteCommentRequest(ctx context.Context, in *postService.CommentIdRequest) (*postService.EmptyRequest, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "DeleteRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	id, _ := primitive.ObjectIDFromHex(in.Id)
	handler.commentService.Delete(ctx, id)
	response := &postService.EmptyRequest{}
	return response, nil
}
