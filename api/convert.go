package api

import (
	db "github.com/ZenSam7/Education/db/sqlc"
	pb "github.com/ZenSam7/Education/protobuf"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convUser(user db.User) *pb.User {
	return &pb.User{
		IdUser:      user.IDUser,
		Name:        user.Name,
		Description: user.Description.String,
		Email:       user.Email,
		Karma:       user.Karma,
		CreatedAt:   timestamppb.New(user.CreatedAt.Time),
	}
}

func convComment(comment db.Comment) *pb.Comment {
	return &pb.Comment{
		IdComment:  comment.IDComment,
		Text:       comment.Text,
		Author:     comment.Author,
		Evaluation: comment.Evaluation,
		CreatedAt:  timestamppb.New(comment.CreatedAt.Time),
		EditedAt:   timestamppb.New(comment.EditedAt.Time),
	}
}

func convArticle(article db.Article) *pb.Article {
	return &pb.Article{
		IdArticle:  article.IDArticle,
		Title:      article.Title,
		Text:       article.Text,
		Comments:   article.Comments,
		Authors:    article.Authors,
		Evaluation: article.Evaluation,
		CreatedAt:  timestamppb.New(article.CreatedAt.Time),
		EditedAt:   timestamppb.New(article.EditedAt.Time),
	}
}

func convImage(image db.Image) *pb.Image {
	return &pb.Image{
		IdImage: image.IDImage,
		Name:    image.Name,
		Content: image.Content,
		IdUser:  image.IDImage,
	}
}
