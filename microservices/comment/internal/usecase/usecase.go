package usecase

import (
	"context"
	"errors"

	"2024_1_kayros/gen/go/comment"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/myerrors/grpcerr"
	"2024_1_kayros/microservices/comment/internal/repo"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

type Comment interface {
	CreateComment(ctx context.Context, comment *comment.Comment) (*comment.Comment, error)
	GetCommentsByRest(ctx context.Context, restId *comment.RestId) (*comment.CommentList, error)
	DeleteComment(ctx context.Context, id *comment.CommentId) (*comment.Empty, error)
}

type CommentLayer struct {
	comment.UnimplementedCommentWorkerServer
	repoComment repo.Comment
	logger      *zap.Logger
}

func NewCommentLayer(repoComProps repo.Comment, loggerProps *zap.Logger) *CommentLayer {
	return &CommentLayer{
		UnimplementedCommentWorkerServer: comment.UnimplementedCommentWorkerServer{},
		repoComment:                      repoComProps,
		logger:                           loggerProps,
	}
}

func (uc *CommentLayer) CreateComment(ctx context.Context, com *comment.Comment) (*comment.Comment, error) {
	res, err := uc.repoComment.Create(ctx, com)
	if err != nil {
		uc.logger.Error(err.Error())
		if errors.Is(err, myerrors.SqlNoRowsCommentRelation) {
			return &comment.Comment{}, grpcerr.NewError(codes.NotFound, myerrors.SqlNoRowsCommentRelation.Error())
		}
		if errors.Is(err, myerrors.SqlNoRowsRestaurantRelation) {
			return &comment.Comment{}, grpcerr.NewError(codes.NotFound, myerrors.SqlNoRowsRestaurantRelation.Error())
		}
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			return &comment.Comment{}, grpcerr.NewError(codes.NotFound, myerrors.SqlNoRowsOrderRelation.Error())
		}
		return &comment.Comment{}, grpcerr.NewError(codes.Internal, err.Error())
	}
	return res, nil
}

func (uc *CommentLayer) GetCommentsByRest(ctx context.Context, restId *comment.RestId) (*comment.CommentList, error) {
	comments, err := uc.repoComment.GetCommentsByRest(ctx, restId)
	if err != nil {
		uc.logger.Error(err.Error())
		return &comment.CommentList{}, grpcerr.NewError(codes.Internal, err.Error())
	}
	return comments, nil
}

func (uc *CommentLayer) DeleteComment(ctx context.Context, id *comment.CommentId) (*comment.Empty, error) {
	err := uc.repoComment.Delete(ctx, id)
	if err != nil {
		uc.logger.Error(err.Error())
		return &comment.Empty{}, grpcerr.NewError(codes.Internal, err.Error())
	}
	return &comment.Empty{}, nil
}
