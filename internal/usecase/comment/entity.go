package comment

import (
	"2024_1_kayros/gen/go/comment"
	"2024_1_kayros/internal/entity"
)

func FromGrpcStructToComment(grpcCom *comment.Comment) *entity.Comment {
	return &entity.Comment{
		Id:        grpcCom.Id,
		UserId:    grpcCom.UserId,
		UserName:  grpcCom.UserName,
		UserImage: grpcCom.Image,
		RestId:    grpcCom.RestId,
		Text:      grpcCom.Text,
		Rating:    grpcCom.Rating,
	}
}

func FromGrpcStructToCommentArray(grpcCom *comment.CommentList) []*entity.Comment {
	if len(grpcCom.GetComment()) == 0 {
		return nil
	}
	comArray := make([]*entity.Comment, len(grpcCom.GetComment()))
	for i, c := range grpcCom.GetComment() {
		comArray[i] = FromGrpcStructToComment(c)
	}
	return comArray
}
