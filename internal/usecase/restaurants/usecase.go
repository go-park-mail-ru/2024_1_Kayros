package rest

import (
	"context"

	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/restaurants"
	"2024_1_kayros/internal/utils/alias"
)

type Usecase interface {
	GetAll(ctx context.Context) ([]*entity.Restaurant, error)
	GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error)
}
type UsecaseLayer struct {
	repoRest restaurants.Repo
	logger   *zap.Logger
}

func NewUsecaseLayer(repoRestProps restaurants.Repo, loggerProps *zap.Logger) Usecase {
	return &UsecaseLayer{
		repoRest: repoRestProps,
		logger:   loggerProps,
	}
}

func (uc *UsecaseLayer) GetAll(ctx context.Context) ([]*entity.Restaurant, error) {
	return uc.repoRest.GetAll(ctx)
}

func (uc *UsecaseLayer) GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error) {
	return uc.repoRest.GetById(ctx, restId)
}
