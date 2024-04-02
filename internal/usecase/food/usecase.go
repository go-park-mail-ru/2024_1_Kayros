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
	AddToOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) (bool, error)
	UpdateCountInOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId, count uint32) (bool, error)
	DeleteFromOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) (bool, error)
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

func (uc *UsecaseLayer) AddToOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) (bool, error) {
	return uc.repoFood.AddToOrder(ctx, foodId, orderId)
}

func (uc *UsecaseLayer) UpdateCountInOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId, count uint32) (bool, error) {
	return uc.repoFood.UpdateCountInOrder(ctx, foodId, orderId, count)
}
func (uc *UsecaseLayer) DeleteFromOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) (bool, error) {
	return uc.repoFood.DeleteFromOrder(ctx, foodId, orderId)
}
