package rest

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/restaurants"
	"2024_1_kayros/internal/utils/alias"
)

type Usecase interface {
	GetAll(ctx context.Context) ([]*entity.Restaurant, error)
	GetById(ctx context.Context, id alias.FoodId) (*entity.Restaurant, error)
}

type UsecaseLayer struct {
	repo *restaurants.Repo
}

func NewUsecaseLayer(r *restaurants.Repo) Usecase {
	return &UsecaseLayer{repo: r}
}

func (uc *UsecaseLayer) GetAll(ctx context.Context) ([]*entity.Restaurant, error) {
	var rests []*entity.Restaurant
	rests, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return rests, err
}

func (uc *UsecaseLayer) GetById(ctx context.Context, id alias.FoodId) (*entity.Restaurant, error) {
	var rest *entity.Restaurant
	rest, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return rest, err
}
