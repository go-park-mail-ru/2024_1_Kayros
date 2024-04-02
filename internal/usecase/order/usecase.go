package order

import (
	"context"
	"time"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/order"
	"2024_1_kayros/internal/repository/user"
)

type UseCaseInterface interface {
	GetBasket(ctx context.Context, email string) (*entity.Order, error)                  // ok
	GetBasketId(ctx context.Context, email string) (uint64, error)                       // ok
	Create(ctx context.Context, email string) error                                      // ok
	Update(ctx context.Context, order *entity.Order) (*entity.Order, error)              // ok
	UpdateStatus(ctx context.Context, orderId uint64, status string) (string, error)     // ok
	AddFoodToOrder(ctx context.Context, foodId int, orderId uint64) error                // ok
	UpdateCountInOrder(ctx context.Context, orderId uint64, foodId int, count int) error // ok
	DeleteFromOrder(ctx context.Context, orderId uint64, foodId int) error               // ok
}

type UseCase struct {
	repoOrder order.RepoInterface
	repoUser  user.Repo
}

func NewUseCase(ro order.RepoInterface, ru user.Repo) UseCaseInterface {
	return &UseCase{repoOrder: ro, repoUser: ru}
}

func (uc *UseCase) GetBasket(ctx context.Context, email string) (*entity.Order, error) {
	status := "Корзина"
	User, err := uc.repoUser.GetByEmail(ctx, email)
	basket, err := uc.repoOrder.GetBasket(ctx, User.Id, status)
	if err != nil {
		return nil, err
	}
	basket.Food, err = uc.repoOrder.GetFood(ctx, basket.Id)
	return basket, nil
}

func (uc *UseCase) GetBasketId(ctx context.Context, email string) (uint64, error) {
	status := "Корзина"
	User, err := uc.repoUser.GetByEmail(ctx, email)
	id, err := uc.repoOrder.GetBasketId(ctx, User.Id, status)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (uc *UseCase) Create(ctx context.Context, email string) error {
	User, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	currentTime := time.Now()
	currentTimeSQL := currentTime.Format("2003-03-03 03:03:03")
	status := "Корзина"
	err = uc.repoOrder.Create(ctx, User.Id, currentTimeSQL, status)
	return err
}

func (uc *UseCase) Update(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	err := uc.repoOrder.Update(ctx, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (uc *UseCase) UpdateStatus(ctx context.Context, orderId uint64, status string) (string, error) {
	return uc.repoOrder.UpdateStatus(ctx, orderId, status)
}

func (uc *UseCase) AddFoodToOrder(ctx context.Context, foodId int, orderId uint64) error {
	count := 1
	return uc.repoOrder.AddToOrder(ctx, orderId, foodId, count)
}

func (uc *UseCase) UpdateCountInOrder(ctx context.Context, orderId uint64, foodId int, count int) error {
	return uc.repoOrder.UpdateCountInOrder(ctx, orderId, foodId, count)
}

func (uc *UseCase) DeleteFromOrder(ctx context.Context, orderId uint64, foodId int) error {
	return uc.repoOrder.DeleteFromOrder(ctx, orderId, foodId)
}
