package food

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/food"
)

type UseCaseInterface interface {
	GetByRest(ctx context.Context, restId uint64) ([]*entity.Food, error)
	GetById(ctx context.Context, id uint64) (*entity.Food, error)
}

type UseCase struct {
	repo food.RepoInterface
}

func NewUseCase(r food.RepoInterface) UseCaseInterface {
	return &UseCase{repo: r}
}

func (uc *UseCase) GetByRest(ctx context.Context, restId uint64) ([]*entity.Food, error) {
	dishes, err := uc.repo.GetByRestId(ctx, restId)
	if err != nil {
		return nil, err
	}
	return dishes, nil
}

func (uc *UseCase) GetById(ctx context.Context, id uint64) (*entity.Food, error) {
	dish, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return dish, nil
}
