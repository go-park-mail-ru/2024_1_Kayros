package rest

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/restaurants"
)

type Usecase interface {
	GetAll(ctx context.Context) ([]*entity.Restaurant, error)
	GetById(ctx context.Context, id int) (*entity.Restaurant, error)
}

type UsecaseLayer struct {
	repo *restaurants.RestaurantRepo
}

func NewUsecase(r *restaurants.RestaurantRepo) Usecase {
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

func (uc *UsecaseLayer) GetById(ctx context.Context, id int) (*entity.Restaurant, error) {
	var rest *entity.Restaurant
	rest, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return rest, err
}
