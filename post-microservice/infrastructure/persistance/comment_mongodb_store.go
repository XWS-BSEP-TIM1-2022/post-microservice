package persistance

import (
	"context"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"post-microservice/model"
)

const (
	COMMENT_COLLECTION = "comments"
)

type CommentMongoDBStore struct {
	comments *mongo.Collection
}

func NewCommentMongoDBStore(client *mongo.Client) model.CommentStore {
	comments := client.Database(DATABASE).Collection(COMMENT_COLLECTION)
	return &CommentMongoDBStore{
		comments: comments,
	}
}

func (store CommentMongoDBStore) Get(ctx context.Context, id primitive.ObjectID) (*model.Comment, error) {
	span := tracer.StartSpanFromContext(ctx, "Get")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	filter := bson.M{"_id": id}
	return store.filterOneComment(ctx, filter)
}

func (store CommentMongoDBStore) GetAll(ctx context.Context) ([]*model.Comment, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.D{{}}
	return store.filterComments(ctx, filter)
}

func (store CommentMongoDBStore) GetAllFromPost(ctx context.Context, postId string) ([]*model.Comment, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllFromPost")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	filter := bson.M{"postId": postId}
	return store.filterComments(ctx, filter)
}

func (store CommentMongoDBStore) Create(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	span := tracer.StartSpanFromContext(ctx, "Create")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	result, err := store.comments.InsertOne(ctx, comment)
	if err != nil {
		return nil, err
	}
	comment.Id = result.InsertedID.(primitive.ObjectID)
	return comment, nil
}

func (store CommentMongoDBStore) Delete(ctx context.Context, id primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "Delete")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	filter := bson.M{"_id": id}
	_, err := store.comments.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (store *CommentMongoDBStore) filterComments(ctx context.Context, filter interface{}) ([]*model.Comment, error) {
	span := tracer.StartSpanFromContext(ctx, "filter")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	cursor, err := store.comments.Find(ctx, filter)
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}
	return decodeComments(ctx, cursor)
}

func (store *CommentMongoDBStore) filterOneComment(ctx context.Context, filter interface{}) (comment *model.Comment, err error) {
	span := tracer.StartSpanFromContext(ctx, "filterOne")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	result := store.comments.FindOne(ctx, filter)
	err = result.Decode(&comment)
	return
}

func decodeComments(ctx context.Context, cursor *mongo.Cursor) (comments []*model.Comment, err error) {
	span := tracer.StartSpanFromContext(ctx, "decode")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	for cursor.Next(ctx) {
		var comment model.Comment
		err = cursor.Decode(&comment)
		if err != nil {
			return
		}
		comments = append(comments, &comment)
	}
	err = cursor.Err()
	return
}
