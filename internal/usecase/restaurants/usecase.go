package rest

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/restaurants"
	"2024_1_kayros/internal/utils/alias"
)

type UseCaseInterface interface {
	GetAll(ctx context.Context) ([]*entity.Restaurant, error)
	GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error)
}

<<<<<<< HEAD
type UsecaseLayer struct {
	repoRest restaurants.Repo
}

func NewUsecaseLayer(repoRestProps restaurants.Repo) Usecase {
	return &UsecaseLayer{repoRest: repoRestProps}
}

func (uc *UsecaseLayer) GetAll(ctx context.Context) ([]*entity.Restaurant, error) {
	return uc.repoRest.GetAll(ctx)
}

func (uc *UsecaseLayer) GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error) {
	return uc.repoRest.GetById(ctx, restId)
=======
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
>>>>>>> 413f5b421db12a295cbeea451991559a66aa908b
}
