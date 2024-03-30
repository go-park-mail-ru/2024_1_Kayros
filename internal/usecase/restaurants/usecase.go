package rest

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/restaurants"
)

type UseCaseInterface interface {
	GetAll(ctx context.Context) ([]*entity.Restaurant, error)
	GetById(ctx context.Context, id int) (*entity.Restaurant, error)
}

type UseCase struct {
	repo restaurants.RepoInterface
}

func NewUseCase(r restaurants.RepoInterface) UseCaseInterface {
	return &UseCase{repo: r}
}

func (uc *UseCase) GetAll(ctx context.Context) ([]*entity.Restaurant, error) {
	rests, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return rests, err
}

func (uc *UseCase) GetById(ctx context.Context, id int) (*entity.Restaurant, error) {
	rest, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return rest, err
}
