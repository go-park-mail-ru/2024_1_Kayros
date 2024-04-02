package order

import (
	"context"
	"time"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/order"
	"2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/utils/alias"
)

type Usecase interface {
	GetBasketByUserEmail(ctx context.Context, email string) (*entity.Order, error)
	GetBasketByUserId(ctx context.Context, userId alias.UserId) (*entity.Order, error)
	GetBasketIdByUserId(ctx context.Context, userId alias.UserId) (alias.OrderId, error)
	GetBasketIdByUserEmail(ctx context.Context, email string) (alias.OrderId, error)
	Create(ctx context.Context, email string) (*entity.Order, error)
	Update(ctx context.Context, order *entity.Order) error
	UpdateStatus(ctx context.Context, orderId alias.OrderId, status string) error
	AddFoodToOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) error
	UpdateCountInOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error
	DeleteFromOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId) error
}

type UsecaseLayer struct {
	repoOrder order.Repo
	repoUser  user.Repo
}

func NewUsecaseLayer(repoOrderProps order.Repo, repoUserProps user.Repo) Usecase {
	return &UsecaseLayer{
		repoOrder: repoOrderProps,
		repoUser:  repoUserProps,
	}
}

func (uc *UsecaseLayer) GetBasketByUserEmail(ctx context.Context, email string) (*entity.Order, error) {
	return uc.repoOrder.GetBasketByUserEmail(ctx, email)
}

func (uc *UsecaseLayer) GetBasketByUserId(ctx context.Context, userId alias.UserId) (*entity.Order, error) {
	return uc.repoOrder.GetBasketByUserId(ctx, userId)
}

func (uc *UsecaseLayer) GetBasketIdByUserId(ctx context.Context, userId alias.UserId) (alias.OrderId, error) {
	return uc.repoOrder.GetBasketIdByUserId(ctx, userId)
}

func (uc *UsecaseLayer) GetBasketIdByUserEmail(ctx context.Context, email string) (alias.OrderId, error) {
	return uc.repoOrder.GetBasketIdByUserEmail(ctx, email)
}

func (uc *UsecaseLayer) Create(ctx context.Context, email string) (*entity.Order, error) {
	u, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	currentTimeSQL := currentTime.Format("2024-04-02 23:34:54")
	ord, err := uc.repoOrder.Create(ctx, alias.UserId(u.Id), currentTimeSQL)
	if err != nil {
		return nil, err
	}
	return ord, err

}

func (uc *UsecaseLayer) Update(ctx context.Context, order *entity.Order) error {
	return uc.repoOrder.Update(ctx, order)
}

func (uc *UsecaseLayer) UpdateStatus(ctx context.Context, orderId alias.OrderId, status string) error {
	return uc.repoOrder.UpdateStatus(ctx, orderId, status)
}

func (uc *UsecaseLayer) AddFoodToOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) error {
	return uc.repoOrder.AddToOrder(ctx, orderId, foodId, 1)
}

func (uc *UsecaseLayer) UpdateCountInOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error {
	return uc.repoOrder.UpdateCountInOrder(ctx, orderId, foodId, count)
}

func (uc *UsecaseLayer) DeleteFromOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId) error {
	return uc.repoOrder.DeleteFromOrder(ctx, orderId, foodId)
}
