package dto

import "2024_1_kayros/internal/entity"

type Comment struct {
	UserId    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	UserImage string `json:"user_img"`
	RestId    int    `json:"rest_id"`
	Text      string `json:"text"`
	Rating    int    `json:"rating"`
}

func NewComment(com *entity.Comment) *Comment {
	return &Comment{
		UserId:    com.UserId,
		UserName:  com.UserName,
		UserImage: com.UserImage,
		RestId:    com.RestId,
		Text:      com.Text,
		Rating:    com.Rating,
	}
}

func NewCommentArray(commentArray []*entity.Comment) []*Comment {
	commentArrayDTO := make([]*Comment, len(commentArray))
	for i, com := range commentArray {
		commentArrayDTO[i] = NewComment(com)
	}
	return commentArrayDTO
}
