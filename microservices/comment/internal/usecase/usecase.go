package usecase

import (
	"context"

	"2024_1_kayros/microservices/comment/internal/repo"
	"2024_1_kayros/microservices/comment/proto"
)

type Comment interface {
	CreateComment(ctx context.Context, comment *comment.Comment) (*comment.Comment, error)
	GetCommentsByRest(ctx context.Context, restId *comment.RestId) (*comment.CommentList, error)
	DeleteComment(ctx context.Context, id *comment.CommentId) (*comment.Empty, error)
}

type CommentLayer struct {
	comment.UnimplementedCommentWorkerServer
	repoComment repo.Comment
}

func NewCommentLayer(repoComProps repo.Comment) *CommentLayer {
	return &CommentLayer{
		UnimplementedCommentWorkerServer: comment.UnimplementedCommentWorkerServer{},
		repoComment:                      repoComProps,
	}
}

func (uc *CommentLayer) CreateComment(ctx context.Context, comment *comment.Comment) (*comment.Comment, error) {
	res, err := uc.repoComment.Create(ctx, comment)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (uc *CommentLayer) GetCommentsByRest(ctx context.Context, restId *comment.RestId) (*comment.CommentList, error) {
	comments, err := uc.repoComment.GetCommentsByRest(ctx, restId)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (uc *CommentLayer) DeleteComment(ctx context.Context, id *comment.CommentId) (*comment.Empty, error) {
	return nil, uc.repoComment.Delete(ctx, id)
}
