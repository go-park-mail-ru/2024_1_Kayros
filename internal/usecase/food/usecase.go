package food

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/food"
)

type UseCaseInterface interface {
	GetByRest(ctx context.Context, restId int) ([]*entity.Food, error)
	GetById(ctx context.Context, id int) (*entity.Food, error)
}

type UseCase struct {
	repo *food.Repo
}

func NewUseCase(r *food.Repo) *UseCase {
	return &UseCase{repo: r}
}

func (uc *UseCase) GetByRest(ctx context.Context, restId int) ([]*entity.Food, error) {
	var dishes []*entity.Food
	dishes, err := uc.repo.GetByRest(ctx, restId)
	if err != nil {
		return nil, err
	}
	return dishes, nil
}

func (uc *UseCase) GetById(ctx context.Context, id int) (*entity.Food, error) {
	var dish *entity.Food
	dish, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return dish, nil
}
