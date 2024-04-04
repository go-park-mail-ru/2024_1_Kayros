package delivery

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"2024_1_kayros/internal/entity/dto"
	ucOrder "2024_1_kayros/internal/usecase/order"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
)

type OrderHandler struct {
	uc     ucOrder.UseCaseInterface
	logger *zap.Logger
}

func NewOrderHandler(u ucOrder.UseCaseInterface, loggerProps *zap.Logger) *OrderHandler {
	return &OrderHandler{
		uc:     u,
		logger: loggerProps,
	}
}

type FoodOrder struct {
	FoodId int `json:"food_id"`
	Count  int `json:"count"`
}

// GET - ok

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := r.Context().Value("email").(string)
	order, err := h.uc.GetBasket(r.Context(), email)
	if err != nil {
		w = functions.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	food := make([]*dto.FoodInOrder, len(order.Food))
	food = dto.NewFoodArray(food, order.Food)
	orderDTO := dto.OrderToDTO(order, food)
	body, err := json.Marshal(orderDTO)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

}

//PUT - ok

func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var order *dto.Order
	email := r.Context().Value("email").(string)
	basket, err := h.uc.GetBasket(r.Context(), email)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusUnauthorized)
		return
	}
	if err = r.Body.Close(); err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusUnauthorized)
		return
	}
	err = json.Unmarshal(body, &order)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	dateOrder, err := time.Parse("2003-03-03 03:03:03", basket.DateOrder)
	dateNew, err := time.Parse("2003-03-03 03:03:03", order.DateReceiving)
	if dateOrder.Before(dateNew) {
		basket.DateReceiving = order.DateReceiving
		basket.Address = order.Address
		basket.ExtraAddress = order.ExtraAddress
	} else {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusUnauthorized)
		return
	}
	basket, err = h.uc.Update(r.Context(), basket)
	if err != nil {
		w = functions.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}
	food := make([]*dto.FoodInOrder, len(basket.Food))
	food = dto.NewFoodArray(food, basket.Food)
	orderDTO := dto.OrderToDTO(basket, food)
	jsonResponse, err := json.Marshal(orderDTO)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// POST-ok

func (h *OrderHandler) AddFood(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	foodId, _ := strconv.Atoi(vars["food_id"])
	email := r.Context().Value("email").(string)
	//передаем почту и статус, чтоб найти id заказа-корзины
	//res uint
	basketId, err := h.uc.GetBasketId(r.Context(), email)
	if err != nil {
		w = functions.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}
	//если res-0, значит у пользователя нет корзины
	if basketId == 0 {
		//создаем заказ-корзину
		err := h.uc.Create(r.Context(), email)
		if err != nil {
			w = functions.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}
	//добавляем еду в заказ
	err = h.uc.AddFoodToOrder(r.Context(), foodId, basketId)
	if err != nil {
		w = functions.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//PUT - ok

func (h *OrderHandler) UpdateFoodCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := r.Context().Value("email").(string)
	basketId, err := h.uc.GetBasketId(r.Context(), email)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusUnauthorized)
		return
	}
	if err = r.Body.Close(); err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusUnauthorized)
		return
	}

	var item FoodOrder
	err = json.Unmarshal(body, &item)

	if err != nil {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	err = h.uc.UpdateCountInOrder(r.Context(), basketId, item.FoodId, item.Count)
	if err != nil {
		w = functions.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//DELETE - ok

func (h *OrderHandler) DeleteFoodFromOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	foodId, _ := strconv.Atoi(vars["food_id"])
	email := r.Context().Value("email").(string)
	basketId, err := h.uc.GetBasketId(r.Context(), email)
	//передаем почту и статус, чтоб найти id заказа-корзины
	//res uint
	err = h.uc.DeleteFromOrder(r.Context(), basketId, foodId)
	if err != nil {
		w = functions.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
