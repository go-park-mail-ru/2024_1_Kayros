package food

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/food"
	"2024_1_kayros/internal/utils/alias"
)

type Usecase interface {
	GetByRest(context.Context, alias.RestId) ([]*entity.Food, error)
	GetById(context.Context, alias.FoodId) (*entity.Food, error)
	AddToOrder(context.Context, alias.FoodId, alias.OrderId) (bool, error)
	UpdateCountInOrder(context.Context, alias.FoodId, alias.OrderId, uint32) (bool, error)
	DeleteFromOrder(context.Context, alias.FoodId, alias.OrderId) (bool, error)
}

type UsecaseLayer struct {
	repo food.Repo
}

func NewUsecaseLayer(r food.Repo) Usecase {
	return &UsecaseLayer{repo: r}
}

func (uc *UsecaseLayer) GetByRest(ctx context.Context, id alias.RestId) ([]*entity.Food, error) {
	var dishes []*entity.Food

	dishes, err := uc.repo.GetByRestId(ctx, id)
	if err != nil {
		return nil, err
	}
	return dishes, nil
}

func (uc *UsecaseLayer) GetById(ctx context.Context, id alias.FoodId) (*entity.Food, error) {
	var dish *entity.Food
	dish, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return dish, nil
}

func (uc *UsecaseLayer) AddToOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) (bool, error) {
	wasAdded, err := uc.repo.AddToOrder(ctx, foodId, orderId)
	return wasAdded, err
}

func (uc *UsecaseLayer) UpdateCountInOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId, count uint32) (bool, error) {
	wasUpdated, err := uc.repo.UpdateCountInOrder(ctx, foodId, orderId, count)
	return wasUpdated, err
}
func (uc *UsecaseLayer) DeleteFromOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) (bool, error) {
	wasDeleted, err := uc.repo.DeleteFromOrder(ctx, foodId, orderId)
	return wasDeleted, err
}
