package comment

import (
	"context"
	"time"

	"2024_1_kayros/gen/go/comment"
	"2024_1_kayros/gen/go/user"
	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/myerrors/grpcerr"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Usecase interface {
	CreateComment(ctx context.Context, comment entity.Comment, email string, orderId uint64) (*entity.Comment, error)
	GetCommentsByRest(ctx context.Context, restId alias.RestId) ([]*entity.Comment, error)
	DeleteComment(ctx context.Context, id uint64) error
}

type UsecaseLayer struct {
	grpcCommentClient comment.CommentWorkerClient
	grpcUserClient    user.UserManagerClient
	metrics           *metrics.Metrics
}

func NewUseCaseLayer(commentClient comment.CommentWorkerClient, userClient user.UserManagerClient, m *metrics.Metrics) Usecase {
	return &UsecaseLayer{
		grpcCommentClient: commentClient,
		grpcUserClient:    userClient,
		metrics:           m,
	}
}

func (uc *UsecaseLayer) CreateComment(ctx context.Context, com entity.Comment, email string, orderId uint64) (*entity.Comment, error) {
	timeNow := time.Now()
	u, err := uc.grpcUserClient.GetData(ctx, &user.Email{Email: email})
	msRequestTimeout := time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.UserMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.CommentMicroservice, grpcStatus.String()).Inc()
		}
		if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
			return &entity.Comment{}, myerrors.SqlNoRowsUserRelation
		}
		return &entity.Comment{}, err
	}
	c := comment.Comment{
		UserId:  u.Id,
		RestId:  com.RestId,
		Text:    com.Text,
		Rating:  com.Rating,
		OrderId: orderId,
	}
	timeNow = time.Now()
	res, err := uc.grpcCommentClient.CreateComment(ctx, &c)
	msRequestTimeout = time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.CommentMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.CommentMicroservice, grpcStatus.String()).Inc()
		}
		if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsCommentRelation) {
			return &entity.Comment{}, myerrors.SqlNoRowsCommentRelation
		}
		if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsRestaurantRelation) {
			return &entity.Comment{}, myerrors.SqlNoRowsRestaurantRelation
		}
		if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsOrderRelation) {
			return &entity.Comment{}, myerrors.SqlNoRowsOrderRelation
		}
	}
	res.UserName = u.Name
	res.Image = u.ImgUrl
	return FromGrpcStructToComment(res), nil
}

// / !!!!!!!
func (uc *UsecaseLayer) GetCommentsByRest(ctx context.Context, restId alias.RestId) ([]*entity.Comment, error) {
	timeNow := time.Now()
	comments, err := uc.grpcCommentClient.GetCommentsByRest(ctx, &comment.RestId{Id: uint64(restId)})
	msRequestTimeout := time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.CommentMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.CommentMicroservice, grpcStatus.String()).Inc()
		}
		return nil, err
	}
	if comments == nil {
		return nil, nil
	}
	return FromGrpcStructToCommentArray(comments), nil
}

// / !!!!!!!!
func (uc *UsecaseLayer) DeleteComment(ctx context.Context, id uint64) error {
	timeNow := time.Now()
	_, err := uc.grpcCommentClient.DeleteComment(ctx, &comment.CommentId{Id: id})
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.CommentMicroservice, grpcStatus.String()).Inc()
		}
	}
	msRequestTimeout := time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.CommentMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	return err
}
