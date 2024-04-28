package rest

import (
	"context"

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
}

func NewUsecaseLayer(repoRestProps restaurants.Repo) Usecase {
	return &UsecaseLayer{
		repoRest: repoRestProps,
	}
}

func (uc *UsecaseLayer) GetAll(ctx context.Context) ([]*entity.Restaurant, error) {
	rests, err := uc.repoRest.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return rests, nil
}

func (uc *UsecaseLayer) GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error) {
	rest, err := uc.repoRest.GetById(ctx, restId)
	if err != nil {
		return nil, err
	}
	return rest, nil
}
