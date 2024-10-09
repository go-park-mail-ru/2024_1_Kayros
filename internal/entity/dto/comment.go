package dto

import "2024_1_kayros/internal/entity"

type Comment struct {
	Id        uint64 `json:"id"`
	UserId    uint64 `json:"user_id,omitempty"`
	UserName  string `json:"user_name"`
	UserImage string `json:"user_img"`
	RestId    uint64 `json:"rest_id,omitempty"`
	Text      string `json:"text"`
	Rating    uint32 `json:"rating"`
}

type CommentArray struct {
	Payload []*Comment `json:"payload" valid:"-"`
}

func NewComment(com *entity.Comment) *Comment {
	return &Comment{
		Id:        com.Id,
		UserId:    com.UserId,
		UserName:  com.UserName,
		UserImage: com.UserImage,
		RestId:    com.RestId,
		Text:      com.Text,
		Rating:    com.Rating,
	}
}

func NewCommentFromDTO(com *Comment) entity.Comment {
	return entity.Comment{
		Id:        com.Id,
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

type InputId struct {
	Id uint64 `json:"id"`
}

type InputComment struct {
	OrderId uint64 `json:"order_id"`
	Text    string `json:"text"`
	Rating  uint32 `json:"rating"`
}
