package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
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

type FoodCount struct {
	FoodId alias.FoodId `json:"food_id" valid:"positive"`
	Count  uint32       `json:"count" valid:"positive"`
}

func ChangingStatus(ctx context.Context, h *OrderHandler, id uint64, arr []string) {
	requestId := ctx.Value(cnst.RequestId)
	for _, s := range arr {
		time.Sleep(10 * time.Second)
		fmt.Println(s)
		_, err := h.uc.UpdateStatus(ctx, alias.OrderId(id), s)
		fmt.Println(err)
		if err != nil {
			h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId.(string)))
			return
		}
	}
}

func (d *FoodCount) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

type payedOrderInfo struct {
	Id     alias.OrderId `json:"id"`
	Status string        `json:"status"`
}

func (h *OrderHandler) GetBasket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)
	if email == "" && unauthId == "" {
		h.logger.Error(myerrors.NoBasket.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusOK)
		return
	}

	var order *entity.Order
	var err error
	if unauthId != "" {
		order, err = h.uc.GetBasketNoAuth(r.Context(), unauthId)
	} else if email != "" {
		order, err = h.uc.GetBasket(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
		} else {
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}
	orderDTO := dto.NewOrder(order)
	w = functions.JsonResponse(w, orderDTO)
}

func (h *OrderHandler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		h.logger.Error(myerrors.AuthorizedEn.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.AuthorizedRu, http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	var order *entity.Order
	order, err = h.uc.GetOrderById(r.Context(), alias.OrderId(id))
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	orderDTO := dto.NewOrder(order)
	w = functions.JsonResponse(w, orderDTO)
}

func (h *OrderHandler) GetCurrentOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		h.logger.Error(myerrors.UnauthorizedEn.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	orders, err := h.uc.GetCurrentOrders(r.Context(), email)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
		} else {
			w = functions.ErrorResponse(w, myerrors.NoOrdersRu, http.StatusInternalServerError)
		}
		return
	}
	ordersDTO := dto.NewShortOrderArray(orders)
	w = functions.JsonResponse(w, ordersDTO)
}

func (h *OrderHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)
	if email == "" && unauthId == "" {
		h.logger.Error(myerrors.NoBasket.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusOK)
		return
	}

	var fullAddress dto.FullAddress
	var basketId alias.OrderId
	var err error
	if unauthId != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), unauthId)
	} else if email != "" {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
		} else {
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	if err = r.Body.Close(); err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &fullAddress)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	isValid, err := fullAddress.Validate()
	if err != nil || !isValid {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	err = h.uc.UpdateAddress(r.Context(), fullAddress, basketId)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	w = functions.JsonResponse(w, map[string]string{"detail": "Адрес заказа был успешно обновлен"})
}

func (h *OrderHandler) Pay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)
	if email == "" && unauthId == "" {
		h.logger.Error(myerrors.AuthorizedEn.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.AuthorizedRu, http.StatusOK)
		return
	}

	var basket *entity.Order
	var err error
	if unauthId != "" {
		basket, err = h.uc.GetBasketNoAuth(r.Context(), unauthId)
	} else if email != "" {
		basket, err = h.uc.GetBasket(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	payedOrder, err := h.uc.Pay(r.Context(), alias.OrderId(basket.Id), basket.Status, email, alias.UserId(basket.UserId))
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.AlreadyPayed) {
			w = functions.ErrorResponse(w, myerrors.AlreadyPayedRu, http.StatusBadRequest)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	statuses := []string{cnst.Created, cnst.Cooking, cnst.OnTheWay, cnst.Delivered}

	ctx := context.Background()
	ctx = context.WithValue(ctx, cnst.RequestId, requestId)
	go ChangingStatus(ctx, h, payedOrder.Id, statuses)

	response := payedOrderInfo{Id: alias.OrderId(payedOrder.Id), Status: payedOrder.Status}
	w = functions.JsonResponse(w, response)

}

func (h *OrderHandler) AddFood(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	if err = r.Body.Close(); err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	var item FoodCount
	err = json.Unmarshal(body, &item)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	isValid, err := item.Validate()
	if err != nil || !isValid {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	var basketId alias.OrderId
	if unauthId != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), unauthId)
	} else if email != "" {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil && !errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
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
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
	}

	// add food in order
	err = h.uc.AddFoodToOrder(r.Context(), item.FoodId, item.Count, basketId)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.OrderAddFood) {
			w = functions.ErrorResponse(w, myerrors.NoAddFoodToOrderRu, http.StatusInternalServerError)
		} else {
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}

	order, err := h.uc.GetOrderById(r.Context(), basketId)
	if err != nil {
		h.logger.Error(myerrors.OrderAddFood.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	orderDTO := dto.NewOrder(order)
	w = functions.JsonResponse(w, orderDTO)
}

func (h *OrderHandler) Clean(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)
	if email == "" && unauthId == "" {
		h.logger.Error(myerrors.NoBasket.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusOK)
		return
	}

	var basketId alias.OrderId
	var err error
	if unauthId != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), unauthId)
	} else if email != "" {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
		} else {
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}

	err = h.uc.Clean(r.Context(), basketId)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.FailCleanBasket) {
			w = functions.ErrorResponse(w, myerrors.FailCleanBasketRu, http.StatusInternalServerError)
		} else {
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}
	w = functions.JsonResponse(w, map[string]string{"detail": "Корзина очищена"})
}

func (h *OrderHandler) UpdateFoodCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)

	var basketId alias.OrderId
	var err error
	if unauthId != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), unauthId)
	} else if email != "" {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
		} else {
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	if err = r.Body.Close(); err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	var item FoodCount
	err = json.Unmarshal(body, &item)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	if item.FoodId <= 0 || item.Count <= 0 {
		h.logger.Error(myerrors.BadCredentialsEn.Error())
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	order, err := h.uc.UpdateCountInOrder(r.Context(), basketId, item.FoodId, item.Count)
	if err != nil {
		h.logger.Error(myerrors.BadCredentialsEn.Error())
		if errors.Is(err, myerrors.SqlNoRowsFoodOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoAddFoodToOrderRu, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	orderDTO := dto.NewOrder(order)
	w = functions.JsonResponse(w, orderDTO)
}

func (h *OrderHandler) DeleteFoodFromOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)
	vars := mux.Vars(r)
	foodId, err := strconv.Atoi(vars["food_id"])
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	var basketId alias.OrderId
	if unauthId != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), unauthId)
	} else if email != "" {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
		} else {
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}
	order, err := h.uc.DeleteFromOrder(r.Context(), basketId, alias.FoodId(foodId))
	if err != nil {
		if errors.Is(err, myerrors.SuccessCleanRu) {
			w = functions.JsonResponse(w, map[string]string{"detail": "Корзина очищена"})
		} else if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoDeleteFoodRu, http.StatusInternalServerError)
		} else {
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		}
		return
	}
	orderDTO := dto.NewOrder(order)
	w = functions.JsonResponse(w, orderDTO)
}
