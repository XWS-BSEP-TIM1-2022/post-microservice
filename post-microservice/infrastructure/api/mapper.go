package api

import (
	postService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/post"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post-microservice/model"
)

func mapPost(post *model.Post) *postService.Post {
	postPb := &postService.Post{
		Id:     post.Id.Hex(),
		UserId: post.UserId,
		Text:   post.Text,
		//Photo: post.Photo,
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
		//Photo: postPb.Photo,
	}
	return post
}
