package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/repository/food"
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
	Pay(ctx context.Context, orderId alias.OrderId, currentStatus string) (*entity.Order, error)
	AddFoodToOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) error
	UpdateCountInOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) (*entity.Order, error)
	DeleteFromOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId) (*entity.Order, error)
}

type UsecaseLayer struct {
	repoOrder order.Repo
	repoUser  user.Repo
	repoFood  food.Repo
	logger    *zap.Logger
}

func NewUsecaseLayer(repoOrderProps order.Repo, repoFoodProps food.Repo, repoUserProps user.Repo, loggerProps *zap.Logger) Usecase {
	return &UsecaseLayer{
		repoOrder: repoOrderProps,
		repoUser:  repoUserProps,
		repoFood:  repoFoodProps,
		logger:    loggerProps,
	}
}

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
	id, err := uc.repoOrder.GetBasketId(ctx, requestId, alias.UserId(u.Id))
	if id == 0 {
		functions.LogWarn(uc.logger, requestId, methodName, err, constants.UsecaseLayer)
		return 0, nil
	}
	if err != nil {
		functions.LogError(uc.logger, requestId, methodName, err, constants.UsecaseLayer)
		return 0, err
	}
	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
	return id, nil
}

func (uc *UsecaseLayer) GetBasket(ctx context.Context, email string) (*entity.Order, error) {
	methodName := constants.NameMethodGetBasket
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	u, err := uc.repoUser.GetByEmail(ctx, email, requestId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	if u == nil {
		functions.LogError(uc.logger, requestId, methodName, err, constants.UsecaseLayer)
		return nil, err
	}
	orders, err := uc.repoOrder.GetOrders(ctx, requestId, alias.UserId(u.Id), constants.Draft)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	if len(orders) == 0 {
		err = errors.New(order.NoBasketError)
		functions.LogInfo(uc.logger, requestId, methodName, order.NoBasketError, constants.UsecaseLayer)
		return nil, err
	}
	basket := orders[0]
	if len(basket.Food) != 0 {
		basket.RestaurantId = basket.Food[0].RestaurantId
	}
	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
	return basket, nil
}

func (uc *UsecaseLayer) Create(ctx context.Context, email string) (alias.OrderId, error) {
	methodName := constants.NameMethodCreateOrder
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	u, err := uc.repoUser.GetByEmail(ctx, email, requestId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return 0, err
	}
	currentTime := time.Now().UTC()
	timeForDB := currentTime.Format("2006-01-02T15:04:05Z07:00")
	id, err := uc.repoOrder.Create(ctx, requestId, alias.UserId(u.Id), timeForDB)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return 0, err
	}
	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
	return id, err
}

func (uc *UsecaseLayer) UpdateAddress(ctx context.Context, FullAddress dto.FullAddress, orderId alias.OrderId) (*entity.Order, error) {
	methodName := constants.NameMethodUpdateAddress
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	id, err := uc.repoOrder.UpdateAddress(ctx, requestId, FullAddress.Address, FullAddress.ExtraAddress, orderId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	updatedOrder, err := uc.repoOrder.GetOrderById(ctx, requestId, id)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	if len(updatedOrder.Food) != 0 {
		updatedOrder.RestaurantId = updatedOrder.Food[0].RestaurantId
	}
	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
	return updatedOrder, nil
}

func (uc *UsecaseLayer) Pay(ctx context.Context, orderId alias.OrderId, currentStatus string) (*entity.Order, error) {
	methodName := constants.NameMethodPayOrder
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	if currentStatus != constants.Draft {
		functions.LogWarn(uc.logger, requestId, methodName, fmt.Errorf("Статус должен быть Draft"), constants.UsecaseLayer)
		return nil, fmt.Errorf("Заказ уже оплачен")
	}
	id, err := uc.repoOrder.UpdateStatus(ctx, requestId, orderId, constants.Payed)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	payedOrder, err := uc.repoOrder.GetOrderById(ctx, requestId, id)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	if len(payedOrder.Food) != 0 {
		payedOrder.RestaurantId = payedOrder.Food[0].RestaurantId
	}
	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
	return payedOrder, nil
}

func (uc *UsecaseLayer) AddFoodToOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) error {
	methodName := constants.NameMethodAddToOrder
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	//inputFood, err := uc.repoFood.GetById(ctx, requestId, foodId)
	//if err != nil {
	//	functions.LogUsecaseFail(uc.logger, requestId, methodName)
	//	return err
	//}
	//fmt.Println(inputFood.RestaurantId)
	//Order, err := uc.repoOrder.GetOrderById(ctx, requestId, orderId)
	//fmt.Println(Order.Food[0].RestaurantId)
	//if inputFood.RestaurantId != Order.Food[0].RestaurantId {
	//	err = uc.repoOrder.CleanBasket(ctx, requestId, orderId)
	//	if err != nil {
	//		functions.LogUsecaseFail(uc.logger, requestId, methodName)
	//		return err
	//	}
	//}
	err := uc.repoOrder.AddToOrder(ctx, requestId, orderId, foodId, 1)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return err
	}
	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
	return err
}

func (uc *UsecaseLayer) UpdateCountInOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) (*entity.Order, error) {
	methodName := constants.NameMethodUpdateCountInOrder
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	err := uc.repoOrder.UpdateCountInOrder(ctx, requestId, orderId, foodId, count)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	updatedOrder, err := uc.repoOrder.GetOrderById(ctx, requestId, orderId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	if len(updatedOrder.Food) != 0 {
		updatedOrder.RestaurantId = updatedOrder.Food[0].RestaurantId
	}
	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
	return updatedOrder, nil
}

func (uc *UsecaseLayer) DeleteFromOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId) (*entity.Order, error) {
	methodName := constants.NameMethodDeleteFromOrder
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	err := uc.repoOrder.DeleteFromOrder(ctx, requestId, orderId, foodId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	updatedOrder, err := uc.repoOrder.GetOrderById(ctx, requestId, orderId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	if len(updatedOrder.Food) != 0 {
		updatedOrder.RestaurantId = updatedOrder.Food[0].RestaurantId
	}
	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
	return updatedOrder, nil
}
