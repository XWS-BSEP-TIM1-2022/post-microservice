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
	DATABASE        = "postsDB"
	POST_COLLECTION = "posts"
)

type PostMongoDBStore struct {
	posts *mongo.Collection
}

func NewPostMongoDBStore(client *mongo.Client) model.PostStore {
	posts := client.Database(DATABASE).Collection(POST_COLLECTION)
	return &PostMongoDBStore{
		posts: posts,
	}
}

func (store *PostMongoDBStore) Get(ctx context.Context, id primitive.ObjectID) (*model.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "Get")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	filter := bson.M{"_id": id}
	return store.filterOne(ctx, filter)
}

func (store *PostMongoDBStore) GetAll(ctx context.Context) ([]*model.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.D{{}}
	return store.filter(ctx, filter)
}

func (store *PostMongoDBStore) GetAllFromUser(ctx context.Context, userId string) ([]*model.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllFromUser")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	filter := bson.M{"userId": userId}
	return store.filter(ctx, filter)
}

func (store *PostMongoDBStore) Create(ctx context.Context, post *model.Post) (*model.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "Create")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	result, err := store.posts.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}
	post.Id = result.InsertedID.(primitive.ObjectID)
	return post, nil
}

func (store *PostMongoDBStore) Update(ctx context.Context, id primitive.ObjectID, post *model.Post) (*model.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "Update")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	updatedPost := bson.M{
		"$set": post,
	}
	filter := bson.M{"_id": id}
	_, err := store.posts.UpdateOne(ctx, filter, updatedPost)

	if err != nil {
		return nil, err
	}
	post.Id = id
	return post, nil
}

func (store *PostMongoDBStore) Delete(ctx context.Context, id primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "Delete")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	filter := bson.M{"_id": id}
	_, err := store.posts.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostMongoDBStore) filter(ctx context.Context, filter interface{}) ([]*model.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "filter")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	cursor, err := store.posts.Find(ctx, filter)
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}
	return decode(ctx, cursor)
}

func (store *PostMongoDBStore) filterOne(ctx context.Context, filter interface{}) (post *model.Post, err error) {
	span := tracer.StartSpanFromContext(ctx, "filterOne")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	result := store.posts.FindOne(ctx, filter)
	err = result.Decode(&post)
	return
}

func decode(ctx context.Context, cursor *mongo.Cursor) (posts []*model.Post, err error) {
	span := tracer.StartSpanFromContext(ctx, "decode")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	for cursor.Next(ctx) {
		var post model.Post
		err = cursor.Decode(&post)
		if err != nil {
			return
		}
		posts = append(posts, &post)
	}
	err = cursor.Err()
	return
}
