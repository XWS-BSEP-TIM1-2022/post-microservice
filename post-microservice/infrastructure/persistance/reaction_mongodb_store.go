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
	REACTION_COLLECTION = "reactions"
)

type ReactionMongoDBStore struct {
	reactions *mongo.Collection
}

func NewReactionMongoDBStore(client *mongo.Client) model.ReactionStore {
	reactions := client.Database(DATABASE).Collection(REACTION_COLLECTION)
	return &ReactionMongoDBStore{
		reactions: reactions,
	}
}

func (store ReactionMongoDBStore) Get(ctx context.Context, id primitive.ObjectID) (*model.Reaction, error) {
	span := tracer.StartSpanFromContext(ctx, "Get")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	filter := bson.M{"_id": id}
	return store.filterOneReaction(ctx, filter)
}

func (store ReactionMongoDBStore) GetAll(ctx context.Context) ([]*model.Reaction, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.D{{}}
	return store.filterReactions(ctx, filter)
}

func (store ReactionMongoDBStore) GetAllFromPost(ctx context.Context, postId string) ([]*model.Reaction, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllFromPost")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	filter := bson.M{"postId": postId}
	return store.filterReactions(ctx, filter)
}

func (store ReactionMongoDBStore) Create(ctx context.Context, comment *model.Reaction) (*model.Reaction, error) {
	span := tracer.StartSpanFromContext(ctx, "Create")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	result, err := store.reactions.InsertOne(ctx, comment)
	if err != nil {
		return nil, err
	}
	comment.Id = result.InsertedID.(primitive.ObjectID)
	return comment, nil
}

func (store ReactionMongoDBStore) Delete(ctx context.Context, id primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "Delete")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	filter := bson.M{"_id": id}
	_, err := store.reactions.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (store *ReactionMongoDBStore) filterReactions(ctx context.Context, filter interface{}) ([]*model.Reaction, error) {
	span := tracer.StartSpanFromContext(ctx, "filter")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	cursor, err := store.reactions.Find(ctx, filter)
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}
	return decodeReactions(ctx, cursor)
}

func (store *ReactionMongoDBStore) filterOneReaction(ctx context.Context, filter interface{}) (reaction *model.Reaction, err error) {
	span := tracer.StartSpanFromContext(ctx, "filterOne")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	result := store.reactions.FindOne(ctx, filter)
	err = result.Decode(&reaction)
	return
}

func decodeReactions(ctx context.Context, cursor *mongo.Cursor) (reactions []*model.Reaction, err error) {
	span := tracer.StartSpanFromContext(ctx, "decode")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(ctx, span)

	for cursor.Next(ctx) {
		var reaction model.Reaction
		err = cursor.Decode(&reaction)
		if err != nil {
			return
		}
		reactions = append(reactions, &reaction)
	}
	err = cursor.Err()
	return
}
