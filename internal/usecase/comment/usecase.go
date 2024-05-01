package comment

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/comment"
	"2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/myerrors"
)

type Usecase interface {
	CreateComment(ctx context.Context, comment entity.Comment, email string) (*entity.Comment, error)
	GetCommentsByRest(ctx context.Context, restId alias.RestId) ([]*entity.Comment, error)
	DeleteComment(ctx context.Context, id uint64) error
}

type UsecaseLayer struct {
	repoComment comment.Repo
	repoUser    user.Repo
}

func NewUseCaseLayer(repoComProps comment.Repo, repoUserProps user.Repo) Usecase {
	return &UsecaseLayer{
		repoComment: repoComProps,
		repoUser:    repoUserProps,
	}
}

func (uc *UsecaseLayer) CreateComment(ctx context.Context, comment entity.Comment, email string) (*entity.Comment, error) {
	u, err := uc.repoUser.GetByEmail(ctx, email)
	comment.UserId = u.Id
	if err != nil {
		return nil, err
	}
	res, err := uc.repoComment.Create(ctx, comment)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (uc *UsecaseLayer) GetCommentsByRest(ctx context.Context, restId alias.RestId) ([]*entity.Comment, error) {
	comments, err := uc.repoComment.GetCommentsByRest(ctx, restId)
	if err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, myerrors.NoComments
	}
	return comments, nil
}

func (uc *UsecaseLayer) DeleteComment(ctx context.Context, id uint64) error {
	return uc.repoComment.Delete(ctx, id)
}
