package delivery

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/entity/dto"
	ucOrder "2024_1_kayros/internal/usecase/order"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
)

type OrderHandler struct {
	uc     ucOrder.Usecase
	logger *zap.Logger
}

func NewOrderHandler(u ucOrder.Usecase, loggerProps *zap.Logger) *OrderHandler {
	return &OrderHandler{
		uc:     u,
		logger: loggerProps,
	}
}

func ChangingStatus(ctx context.Context, h *OrderHandler, id uint64, arr []string) {
	requestId := ctx.Value(cnst.RequestId)
	for _, s := range arr {
		time.Sleep(10 * time.Second)
		_, err := h.uc.UpdateStatus(ctx, alias.OrderId(id), s)
		if err != nil {
			h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId.(string)))
			return
		}
	}
}

func (h *OrderHandler) GetBasket(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)
	if email == "" && unauthId == "" {
		h.logger.Error(myerrors.NoBasket.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusOK)
		return
	}

	var order *entity.Order
	var err error
	if unauthId != "" {
		order, err = h.uc.GetBasketNoAuth(r.Context(), unauthId)
	}
	if email != "" && order == nil {
		order, err = h.uc.GetBasket(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
		} else {
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}

	//проверка промокода

	//есть basketId
	//надо чекнуть наличие промокода в заказе
	promocode, err := h.uc.GetPromocodeByOrder(r.Context(), (*alias.OrderId)(&order.Id))
	//fmt.Println(promocode, err)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	//есть промокод
	if promocode != nil && promocode.Id != 0 && email != "" {
		//но нет кук, надо удалить промик
		if email == "" {
			err = h.uc.DeletePromocode(r.Context(), alias.OrderId(order.Id))
			if err != nil {
				h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			}
			order.Error = cnst.PromocodeIsDeleted
		} else {
			//куки есть, надо проверить актуальность промокода
			_, err = h.uc.CheckPromocode(r.Context(), email, promocode.Code, alias.OrderId(order.Id))
			//если неактуален, то удалить
			if errors.Is(err, myerrors.OverDatePromocode) || errors.Is(err, myerrors.OncePromocode) || errors.Is(err, myerrors.SumPromocode) {
				err = h.uc.DeletePromocode(r.Context(), alias.OrderId(order.Id))
				promocode = nil
				if err != nil {
					h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				}
				order.Error = cnst.PromocodeIsDeleted
			}
		}
	}

	if promocode != nil {
		order.Promocode = promocode.Code
		order.NewSum = order.Sum * (100 - uint64(promocode.Sale)) / 100
	}

	orderDTO := dto.NewOrder(order)
	functions.JsonResponse(w, orderDTO)
}

func (h *OrderHandler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		h.logger.Error(myerrors.AuthorizedEn.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.AuthorizedRu, http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	var order *entity.Order
	order, err = h.uc.GetOrderById(r.Context(), alias.OrderId(id))
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	orderDTO := dto.NewOrder(order)
	functions.JsonResponse(w, orderDTO)
}

func (h *OrderHandler) GetCurrentOrders(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		h.logger.Error(myerrors.UnauthorizedEn.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	orders, err := h.uc.GetCurrentOrders(r.Context(), email)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
		} else {
			functions.ErrorResponse(w, myerrors.NoOrdersRu, http.StatusInternalServerError)
		}
		return
	}
	ordersDtoArray := &dto.ShortOrderArray{Payload: dto.NewShortOrderArray(orders)}
	functions.JsonResponse(w, ordersDtoArray)
}

func (h *OrderHandler) GetArchiveOrders(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		h.logger.Error(myerrors.UnauthorizedEn.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	orders, err := h.uc.GetArchiveOrders(r.Context(), email)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.NoOrdersRu, http.StatusInternalServerError)
		return
	}
	ordersDtoArray := &dto.ShortOrderArray{Payload: dto.NewShortOrderArray(orders)}
	functions.JsonResponse(w, ordersDtoArray)
}

func (h *OrderHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)
	if email == "" && unauthId == "" {
		h.logger.Error(myerrors.NoBasket.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusOK)
		return
	}

	var fullAddress dto.FullAddress
	var basketId alias.OrderId
	var err error
	if unauthId != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), unauthId)
	}
	if email != "" && basketId == 0 {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
		} else {
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	if err = r.Body.Close(); err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	err = easyjson.Unmarshal(body, &fullAddress)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	isValid, err := fullAddress.Validate()
	if err != nil || !isValid {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	err = h.uc.UpdateAddress(r.Context(), fullAddress, basketId)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	functions.JsonResponse(w, &dto.ResponseDetail{Detail: "Адрес заказа был успешно обновлен"})
}

// добавить проверку промокодов
func (h *OrderHandler) Pay(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)
	if email == "" && unauthId == "" {
		h.logger.Error(myerrors.AuthorizedEn.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.AuthorizedRu, http.StatusUnauthorized)
		return
	}

	var basket *entity.Order
	var err error
	if unauthId != "" {
		basket, err = h.uc.GetBasketNoAuth(r.Context(), unauthId)
	}
	if email != "" && basket == nil {
		basket, err = h.uc.GetBasket(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
			return
		}
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	//надо чекнуть наличие промокода в заказе
	promocode, err := h.uc.GetPromocodeByOrder(r.Context(), (*alias.OrderId)(&basket.Id))
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	//есть промокод
	if promocode != nil && promocode.Id != 0 {
		//надо проверить актуальность промокода
		_, err = h.uc.CheckPromocode(r.Context(), email, promocode.Code, alias.OrderId(basket.Id))
		//если неактуален, то удалить
		if errors.Is(err, myerrors.OverDatePromocode) || errors.Is(err, myerrors.OncePromocode) || errors.Is(err, myerrors.SumPromocode) {
			err = h.uc.DeletePromocode(r.Context(), alias.OrderId(basket.Id))
			if err != nil {
				h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			}
			functions.JsonResponse(w, &dto.ResponseDetail{Detail: cnst.PromocodeIsDeleted})
			return
		} else if err != nil {
			h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		} else {
			sum := basket.Sum * (100 - uint64(promocode.Sale)) / 100
			err = h.uc.UpdateSum(r.Context(), sum, alias.OrderId(basket.Id))
			if err != nil {
				h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			}
		}
	}

	payedOrder, err := h.uc.Pay(r.Context(), alias.OrderId(basket.Id), basket.Status, email, alias.UserId(basket.UserId))
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.AlreadyPayed) {
			functions.ErrorResponse(w, myerrors.AlreadyPayedRu, http.StatusBadRequest)
			return
		}
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	statuses := []string{cnst.Cooking, cnst.OnTheWay, cnst.Delivered}

	ctx := context.WithValue(context.Background(), cnst.RequestId, requestId)
	go ChangingStatus(ctx, h, payedOrder.Id, statuses)
	response := &dto.PayedOrderInfo{Id: alias.OrderId(payedOrder.Id), Status: payedOrder.Status}
	functions.JsonResponse(w, response)
}

// добавить проверку промокодов - добавила
func (h *OrderHandler) AddFood(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)
	if email == "" && unauthId == "" {
		h.logger.Error(myerrors.AuthorizedEn.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.AuthorizedRu, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	if err = r.Body.Close(); err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	var item dto.FoodCount
	err = easyjson.Unmarshal(body, &item)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	isValid, err := item.Validate()
	if err != nil || !isValid {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	var basketId alias.OrderId
	if unauthId != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), unauthId)
	}
	if email != "" && basketId == 0 {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil && !errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	if basketId == 0 {
		if email != "" {
			basketId, err = h.uc.Create(r.Context(), email)
		} else if unauthId != "" {
			basketId, err = h.uc.CreateNoAuth(r.Context(), unauthId)
		}
		if err != nil {
			h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
	}

	// add food in order
	err = h.uc.AddFoodToOrder(r.Context(), item.FoodId, item.Count, basketId)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.OrderAddFood) {
			functions.ErrorResponse(w, myerrors.NoAddFoodToOrderRu, http.StatusInternalServerError)
		} else {
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}

	//проверка промокода

	//есть basketId
	//надо чекнуть наличие промокода в заказе
	promocode, err := h.uc.GetPromocodeByOrder(r.Context(), &basketId)
	//fmt.Println(promocode, err)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	//есть промокод
	if promocode != nil && promocode.Id != 0 {
		//но нет кук, надо удалить промик
		if email == "" {
			err = h.uc.DeletePromocode(r.Context(), basketId)
			if err != nil {
				h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			}
		} else {
			//куки есть, надо проверить актуальность промокода
			_, err = h.uc.CheckPromocode(r.Context(), email, promocode.Code, basketId)
			//если неактуален, то удалить
			if errors.Is(err, myerrors.OverDatePromocode) || errors.Is(err, myerrors.OncePromocode) || errors.Is(err, myerrors.SumPromocode) {
				err = h.uc.DeletePromocode(r.Context(), basketId)
				promocode = nil
				if err != nil {
					h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				}
			}
		}
	}

	order, err := h.uc.GetOrderById(r.Context(), basketId)
	if err != nil {
		h.logger.Error(myerrors.OrderAddFood.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	if promocode != nil {
		order.Promocode = promocode.Code
		order.NewSum = order.Sum * (100 - uint64(promocode.Sale)) / 100
	}

	orderDTO := dto.NewOrder(order)
	functions.JsonResponse(w, orderDTO)
}

func (h *OrderHandler) Clean(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)
	if email == "" && unauthId == "" {
		h.logger.Error(myerrors.NoBasket.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusOK)
		return
	}

	var basketId alias.OrderId
	var err error
	if unauthId != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), unauthId)
	}
	if email != "" && basketId == 0 {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
		} else {
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}

	err = h.uc.Clean(r.Context(), basketId)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.FailCleanBasket) {
			functions.ErrorResponse(w, myerrors.FailCleanBasketRu, http.StatusInternalServerError)
		} else {
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}
	functions.JsonResponse(w, &dto.ResponseDetail{Detail: "Корзина очищена"})
}

// добавить проверку промокодов
func (h *OrderHandler) UpdateFoodCount(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)

	var basketId alias.OrderId
	var err error
	if unauthId != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), unauthId)
	}
	if email != "" && basketId == 0 {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
		} else {
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	if err = r.Body.Close(); err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	var item dto.FoodCount
	err = easyjson.Unmarshal(body, &item)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	if item.FoodId <= 0 || item.Count <= 0 {
		h.logger.Error(myerrors.BadCredentialsEn.Error())
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	order, err := h.uc.UpdateCountInOrder(r.Context(), basketId, item.FoodId, item.Count)
	if err != nil {
		h.logger.Error(myerrors.BadCredentialsEn.Error())
		if errors.Is(err, myerrors.SqlNoRowsFoodOrderRelation) {
			functions.ErrorResponse(w, myerrors.NoAddFoodToOrderRu, http.StatusInternalServerError)
			return
		}
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	//проверка промокода

	//есть basketId
	//надо чекнуть наличие промокода в заказе
	promocode, err := h.uc.GetPromocodeByOrder(r.Context(), &basketId)
	//fmt.Println(promocode, err)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	//есть промокод
	if promocode != nil && promocode.Id != 0 {
		//но нет кук, надо удалить промик
		if email == "" {
			err = h.uc.DeletePromocode(r.Context(), basketId)
			if err != nil {
				h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			}
			order.Error = cnst.PromocodeIsDeleted
		} else {
			//куки есть, надо проверить актуальность промокода
			_, err = h.uc.CheckPromocode(r.Context(), email, promocode.Code, basketId)
			//если неактуален, то удалить
			if errors.Is(err, myerrors.OverDatePromocode) || errors.Is(err, myerrors.OncePromocode) || errors.Is(err, myerrors.SumPromocode) {
				err = h.uc.DeletePromocode(r.Context(), basketId)
				promocode = nil
				if err != nil {
					h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				}
				order.Error = cnst.PromocodeIsDeleted
			}
		}
	}

	if promocode != nil {
		order.Promocode = promocode.Code
		order.NewSum = order.Sum * (100 - uint64(promocode.Sale)) / 100
	}

	orderDTO := dto.NewOrder(order)
	functions.JsonResponse(w, orderDTO)
}

// добавить проверку промокодов
func (h *OrderHandler) DeleteFoodFromOrder(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)
	vars := mux.Vars(r)
	foodId, err := strconv.Atoi(vars["food_id"])
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	var basketId alias.OrderId
	if unauthId != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), unauthId)
	}
	if email != "" && basketId == 0 {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
		} else {
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}
	order, err := h.uc.DeleteFromOrder(r.Context(), basketId, alias.FoodId(foodId))
	if err != nil {
		if errors.Is(err, myerrors.SuccessCleanRu) {
			functions.JsonResponse(w, &dto.ResponseDetail{Detail: "Корзина очищена"})
		} else if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			functions.ErrorResponse(w, myerrors.NoDeleteFoodRu, http.StatusInternalServerError)
		} else {
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}

	//проверка промокода

	//есть basketId
	//надо чекнуть наличие промокода в заказе
	promocode, err := h.uc.GetPromocodeByOrder(r.Context(), &basketId)
	//fmt.Println(promocode, err)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	//есть промокод
	if promocode != nil && promocode.Id != 0 {
		//но нет кук, надо удалить промик
		if email == "" {
			err = h.uc.DeletePromocode(r.Context(), basketId)
			if err != nil {
				h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			}
			order.Error = cnst.PromocodeIsDeleted
		} else {
			//куки есть, надо проверить актуальность промокода
			_, err = h.uc.CheckPromocode(r.Context(), email, promocode.Code, basketId)
			//если неактуален, то удалить
			if errors.Is(err, myerrors.OverDatePromocode) || errors.Is(err, myerrors.OncePromocode) || errors.Is(err, myerrors.SumPromocode) {
				err = h.uc.DeletePromocode(r.Context(), basketId)
				promocode = nil
				if err != nil {
					h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				}
				order.Error = cnst.PromocodeIsDeleted
			}
		}
	}

	if promocode != nil {
		order.Promocode = promocode.Code
		order.NewSum = order.Sum * (100 - uint64(promocode.Sale)) / 100
	}

	orderDTO := dto.NewOrder(order)
	functions.JsonResponse(w, orderDTO)
}

func (h *OrderHandler) SetPromocode(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)
	if email == "" && unauthId == "" {
		h.logger.Error(myerrors.NoBasket.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusOK)
		return
	}
	if email == "" {
		h.logger.Error(myerrors.AuthorizedEn.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.AuthorizedRu, http.StatusUnauthorized)
		return
	}

	//reading body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	if err = r.Body.Close(); err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	c := dto.Promo{}
	err = easyjson.Unmarshal(body, &c)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	code := c.Code

	var basketId alias.OrderId
	//берем корзину, к которой может быть применен
	if unauthId != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), unauthId)
	}
	if basketId == 0 {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
			return
		}
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	//checking if promocode valid, actual and etc.
	codeInfo, err := h.uc.CheckPromocode(r.Context(), email, code, basketId)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.OverDatePromocode) {
			functions.ErrorResponse(w, myerrors.OverDatePromocodeRu, http.StatusOK)
			return
		}
		if errors.Is(err, myerrors.OncePromocode) {
			functions.ErrorResponse(w, myerrors.OncePromocodeRu, http.StatusOK)
			return
		}
		if errors.Is(err, myerrors.SumPromocode) {
			functions.ErrorResponse(w, myerrors.SumPromocodeRu, http.StatusOK)
			return
		}
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	if codeInfo == nil {
		functions.JsonResponse(w, &dto.ResponseDetail{Detail: "Такого промокода нет"})
		return
	}

	//фух, вроде все проверила, промокод может быть применен

	//now can set promocode to basket
	saleSum, err := h.uc.SetPromocode(r.Context(), basketId, codeInfo)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsPromocodeRelation) {
			functions.ErrorResponse(w, myerrors.NoSetPromocodeRu, http.StatusOK)
		} else {
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}
	codeDTO := &dto.Promo{
		Id:     codeInfo.Id,
		Code:   codeInfo.Code,
		NewSum: saleSum,
	}
	functions.JsonResponse(w, codeDTO)
}

func (h *OrderHandler) GetAllPromocode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	codes, err := h.uc.GetAllPromocode(r.Context())
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.NoOrdersRu, http.StatusInternalServerError)
		return
	}
	codesDtoArray := &dto.PromocodeArray{Payload: dto.NewPromocodeArray(codes)}
	functions.JsonResponse(w, codesDtoArray)
}
