package order

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/repository/order"
	"2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
)

type Usecase interface {
	GetBasketId(ctx context.Context, email string) (alias.OrderId, error)
	GetBasket(ctx context.Context, email string) (*entity.Order, error)
	Create(ctx context.Context, email string) (alias.OrderId, error)
	UpdateAddress(ctx context.Context, FullAddress dto.FullAddress, orderId alias.OrderId) (*entity.Order, error)
	Pay(ctx context.Context, requestId string, orderId alias.OrderId, currentStatus string) (*entity.Order, error)
	AddFoodToOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) error
	UpdateCountInOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error
	DeleteFromOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId) error
}

type UsecaseLayer struct {
	repoOrder order.Repo
	repoUser  user.Repo
	logger    *zap.Logger
}

func NewUsecaseLayer(repoOrderProps order.Repo, repoUserProps user.Repo, loggerProps *zap.Logger) Usecase {
	return &UsecaseLayer{
		repoOrder: repoOrderProps,
		repoUser:  repoUserProps,
		logger:    loggerProps,
	}
}

// ok
func (uc *UsecaseLayer) GetBasketId(ctx context.Context, email string) (alias.OrderId, error) {
	methodName := constants.NameMethodGetBasketId
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	u, err := uc.repoUser.GetByEmail(ctx, email, requestId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return 0, err
	}
	if u == nil {
		functions.LogError(uc.logger, requestId, methodName, fmt.Errorf("No user"), constants.UsecaseLayer)
		return 0, err
	}
	id, err := uc.repoOrder.GetBasketId(ctx, alias.UserId(u.Id))
	if err != nil {
		functions.LogError(uc.logger, requestId, methodName, err, constants.UsecaseLayer)
		return 0, err
	}
	if id == 0 {
		return 0, fmt.Errorf("Корзина пуста")
	}
	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
	return id, nil
}

// ok
func (uc *UsecaseLayer) GetBasket(ctx context.Context, email string) (*entity.Order, error) {
	fmt.Println("uc order")
	methodName := constants.NameMethodGetBasket
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	fmt.Println("uc getbasket", email)
	u, err := uc.repoUser.GetByEmail(ctx, email, requestId)
	fmt.Println("user_id, status", u.Id, constants.Draft)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	if u == nil {
		functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
		return nil, err
	}
	orders, err := uc.repoOrder.GetOrders(ctx, alias.UserId(u.Id), constants.Draft)
	if err != nil {
		fmt.Println("uc, err with orders: ", err)
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	if len(orders) == 0 {
		functions.LogError(uc.logger, requestId, methodName, fmt.Errorf(order.NoBasketError), constants.UsecaseLayer)
		return nil, fmt.Errorf(order.NoBasketError)
	}
	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
	return orders[0], nil
}

// ok
func (uc *UsecaseLayer) Create(ctx context.Context, email string) (alias.OrderId, error) {
	methodName := constants.NameMethodCreateOrder
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	u, err := uc.repoUser.GetByEmail(ctx, email, requestId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return 0, err
	}

	dateOrder := time.Now().UTC()
	dateOrderForDB := dateOrder.Format("2006-01-02 15:04:05-07:00")
	id, err := uc.repoOrder.Create(ctx, alias.UserId(u.Id), dateOrderForDB)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return 0, err
	}
	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
	return id, err
}

func (uc *UsecaseLayer) UpdateAddress(ctx context.Context, FullAddress dto.FullAddress, orderId alias.OrderId) (*entity.Order, error) {
	methodName := constants.NameMethodUpdateOrder
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	id, err := uc.repoOrder.UpdateAddress(ctx, FullAddress.Address, FullAddress.ExtraAddress, orderId)
	fmt.Println(id, err)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	updatedOrder, err := uc.repoOrder.GetOrderById(ctx, id)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
	return updatedOrder, nil
}

func (uc *UsecaseLayer) Pay(ctx context.Context, requestId string, orderId alias.OrderId, currentStatus string) (*entity.Order, error) {
	if currentStatus != constants.Draft {
		return nil, fmt.Errorf("Заказ уже оплачен")
	}
	id, err := uc.repoOrder.UpdateStatus(ctx, orderId, constants.Payed)
	fmt.Println(id, err)
	if err != nil {
		functions.LogError(uc.logger, requestId, constants.NamePayOrder, err, constants.UsecaseLayer)
		return nil, err
	}
	payedOrder, err := uc.repoOrder.GetOrderById(ctx, id)
	fmt.Println(payedOrder.Id, payedOrder.Status, err)
	if err != nil {
		functions.LogError(uc.logger, requestId, constants.NamePayOrder, err, constants.UsecaseLayer)
		return nil, err
	}
	functions.LogOk(uc.logger, requestId, constants.NamePayOrder, constants.UsecaseLayer)
	return payedOrder, nil
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
