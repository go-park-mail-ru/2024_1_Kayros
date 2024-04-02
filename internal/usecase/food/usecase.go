package food

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/food"
	"2024_1_kayros/internal/utils/alias"
)

type Usecase interface {
	GetByRest(ctx context.Context, restId alias.RestId) ([]*entity.Food, error)
	GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error)
}

type UsecaseLayer struct {
	repoFood food.Repo
}

func NewUsecaseLayer(repoFoodProps food.Repo) Usecase {
	return &UsecaseLayer{repoFood: repoFoodProps}
}

func (uc *UsecaseLayer) GetByRest(ctx context.Context, restId alias.RestId) ([]*entity.Food, error) {
	return uc.repoFood.GetByRestId(ctx, restId)
}

func (uc *UsecaseLayer) GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error) {
	return uc.repoFood.GetById(ctx, foodId)
}
