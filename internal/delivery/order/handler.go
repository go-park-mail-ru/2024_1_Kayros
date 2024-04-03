package delivery

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	orderStatus "2024_1_kayros/internal/utils/constants"
	"github.com/gorilla/mux"

	"2024_1_kayros/internal/entity/dto"
	ucOrder "2024_1_kayros/internal/usecase/order"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
)

type Delivery struct {
	ucOrder ucOrder.Usecase
}

func NewDeliveryLayer(ucOrderProps ucOrder.Usecase) *Delivery {
	return &Delivery{ucOrder: ucOrderProps}
}

func (d *Delivery) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := r.Context().Value("email")
	if email == nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		status = orderStatus.Draft
	}
	basket, err := d.ucOrder.GetOrdersByUserEmail(r.Context(), email.(string), status)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.NotFoundError, http.StatusNotFound)
		return
	}

	basketDTO := dto.NewOrders(basket)
	w = functions.JsonResponse(w, basketDTO)
}

// PUT - ok
func (d *Delivery) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := r.Context().Value("email")
	if email == nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	var orderDTO *dto.Order
	err = json.Unmarshal(body, &orderDTO)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}

	isValid, err := orderDTO.Validate()
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if !isValid {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}

	order := dto.NewOrderFromDTO(orderDTO)
	err = d.ucOrder.Update(r.Context(), order)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	w = functions.JsonResponse(w, "Данные успешно обновлены")
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
