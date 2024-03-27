package food

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/food"
)

type UseCaseInterface interface {
	GetByRest(ctx context.Context, restId int) ([]*entity.Food, error)
	GetById(ctx context.Context, id int) (*entity.Food, error)
	AddToOrder(ctx context.Context, foodId int, orderId int) error
	UpdateCountInOrder(ctx context.Context, foodId int, orderId int, count int) error
	DeleteFromOrder(ctx context.Context, foodId int, orderId int) error
}

type UseCase struct {
	repo food.Repo
}

func NewUseCase(r food.Repo) *UseCase {
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

func (uc *UseCase) AddToOrder(ctx context.Context, foodId int, orderId int) error {
	err := uc.repo.AddToOrder(ctx, foodId, orderId)
	return err
}

func (uc *UseCase) UpdateCountInOrder(ctx context.Context, foodId int, orderId int, count int) error {
	err := uc.repo.UpdateCountInOrder(ctx, foodId, orderId, count)
	return err
}
func (uc *UseCase) DeleteFromOrder(ctx context.Context, foodId int, orderId int) error {
	err := uc.repo.DeleteFromOrder(ctx, foodId, orderId)
	return err
}
