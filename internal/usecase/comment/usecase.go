package comment

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/utils/alias"
	comment "2024_1_kayros/microservices/comment/proto"
)

type Usecase interface {
	CreateComment(ctx context.Context, comment entity.Comment, email string) (*entity.Comment, error)
	GetCommentsByRest(ctx context.Context, restId alias.RestId) ([]*entity.Comment, error)
	DeleteComment(ctx context.Context, id uint64) error
}

type UsecaseLayer struct {
	grpcClient comment.CommentWorkerClient
	repoUser   user.Repo
}

func NewUseCaseLayer(commentClient comment.CommentWorkerClient, repoUserProps user.Repo) Usecase {
	return &UsecaseLayer{
		grpcClient: commentClient,
		repoUser:   repoUserProps,
	}
}

func (uc *UsecaseLayer) CreateComment(ctx context.Context, com entity.Comment, email string) (*entity.Comment, error) {
	u, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	c := comment.Comment{
		UserId: u.Id,
		RestId: com.RestId,
		Text:   com.Text,
		Rating: com.Rating,
	}
	res, err := uc.grpcClient.CreateComment(ctx, &c)
	if err != nil {
		return nil, err
	}
	res.UserName = u.Name
	res.Image = u.ImgUrl
	return FromGrpcStructToComment(res), nil
}

func (uc *UsecaseLayer) GetCommentsByRest(ctx context.Context, restId alias.RestId) ([]*entity.Comment, error) {
	comments, err := uc.grpcClient.GetCommentsByRest(ctx, &comment.RestId{Id: uint64(restId)})
	if err != nil {
		return nil, err
	}
	if comments == nil {
		return nil, nil
	}
	return FromGrpcStructToCommentArray(comments), nil
}

func (uc *UsecaseLayer) DeleteComment(ctx context.Context, id uint64) error {
	_, err := uc.grpcClient.DeleteComment(ctx, &comment.CommentId{Id: id})
	return err
}
