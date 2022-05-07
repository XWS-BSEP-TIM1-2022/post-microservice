package api

import (
	postService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/post"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post-microservice/model"
)

func mapPost(post *model.Post) *postService.Post {
	postPb := &postService.Post{
		Id:           post.Id.Hex(),
		UserId:       post.UserId,
		Text:         post.Text,
		Image:        post.Image,
		Links:        post.Links,
		CreationDate: post.CreationDate.String(),
	}
	return postPb
}

func mapPostPb(postPb *postService.Post) *model.Post {
	id, _ := primitive.ObjectIDFromHex(postPb.Id)

	post := &model.Post{
		Id:     id,
		UserId: postPb.UserId,
		Text:   postPb.Text,
		Links:  postPb.Links,
		Image:  postPb.Image,
	}
	return post
}

func mapComment(comment *model.Comment) *postService.Comment {
	commentPb := &postService.Comment{
		Id:           comment.Id.Hex(),
		UserId:       comment.UserId,
		PostId:       comment.PostId,
		Text:         comment.Text,
		CreationDate: comment.CreationDate.String(),
	}
	return commentPb
}

func mapCommentPb(commentPb *postService.Comment) *model.Comment {
	id, _ := primitive.ObjectIDFromHex(commentPb.Id)

	comment := &model.Comment{
		Id:     id,
		UserId: commentPb.UserId,
		PostId: commentPb.PostId,
		Text:   commentPb.Text,
	}
	return comment
}

func mapReaction(reaction *model.Reaction) *postService.Reaction {
	reactionPb := &postService.Reaction{
		Id:     reaction.Id.Hex(),
		UserId: reaction.UserId,
		PostId: reaction.PostId,
		Type:   reaction.Type,
	}
	return reactionPb
}

func mapReactionPb(reactionPb *postService.Reaction) *model.Reaction {
	id, _ := primitive.ObjectIDFromHex(reactionPb.Id)

	reaction := &model.Reaction{
		Id:     id,
		UserId: reactionPb.UserId,
		PostId: reactionPb.PostId,
		Type:   reactionPb.Type,
	}
	return reaction
}
