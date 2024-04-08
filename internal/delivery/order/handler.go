package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"2024_1_kayros/internal/entity/dto"
	repoErrors "2024_1_kayros/internal/repository/order"
	ucOrder "2024_1_kayros/internal/usecase/order"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/constants"
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
	FoodId int `json:"food_id"`
	Count  int `json:"count"`
}

// GET - ok

func (h *OrderHandler) GetBasket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		err := errors.New("request_id передан не был")
		functions.LogError(h.logger, requestId, constants.NameHandlerSignUp, err, constants.DeliveryLayer)
	} else {
		requestId = ctxRequestId.(string)
	}
	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
		fmt.Println(ctxEmail.(string))
	}
	if email == "" {
		functions.LogErrorResponse(h.logger, requestId, constants.NameMethodGetBasket, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}
	order, err := h.uc.GetBasket(r.Context(), email)
	if err != nil {
		if err.Error() == repoErrors.NoBasketError {
			functions.LogErrorResponse(h.logger, requestId, constants.NameHandlerUserData, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, constants.DeliveryLayer)
			w = functions.ErrorResponse(w, repoErrors.NoBasketError, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	orderDTO := dto.NewOrder(order)
	body, err := json.Marshal(orderDTO)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//PUT - ok

func (h *OrderHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var fullAddress dto.FullAddress
	email := r.Context().Value("email").(string)
	basketId, err := h.uc.GetBasketId(r.Context(), email)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusUnauthorized)
		return
	}
	if err = r.Body.Close(); err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusUnauthorized)
		return
	}
	err = json.Unmarshal(body, &fullAddress)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	basket, err := h.uc.UpdateAddress(r.Context(), fullAddress, basketId)
	if err != nil {
		if err.Error() == repoErrors.NotUpdateError {
			w = functions.ErrorResponse(w, repoErrors.NotUpdateError, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	orderDTO := dto.NewOrder(basket)
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
	fmt.Println("ok")
}

//PUT - ok

func (h *OrderHandler) Pay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		err := errors.New("request_id передан не был")
		functions.LogError(h.logger, requestId, constants.NameHandlerSignUp, err, constants.DeliveryLayer)
	} else {
		requestId = ctxRequestId.(string)
	}
	email := r.Context().Value("email").(string)
	basket, err := h.uc.GetBasket(r.Context(), email)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NamePayOrder, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	payedOrder, err := h.uc.Pay(r.Context(), requestId, alias.OrderId(basket.Id), basket.Status)
	if err != nil {
		if err.Error() == repoErrors.NotUpdateStatusError {
			w = functions.ErrorResponse(w, repoErrors.NotUpdateStatusError, http.StatusInternalServerError)
			return
		}
		functions.LogError(h.logger, requestId, constants.NamePayOrder, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusUnauthorized)
		return
	}
	orderDTO := dto.NewOrder(payedOrder)
	jsonResponse, err := json.Marshal(orderDTO)
	if err != nil {
		functions.LogUsecaseFail(h.logger, requestId, constants.NamePayOrder)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NamePayOrder, err, constants.UsecaseLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	functions.LogOk(h.logger, requestId, constants.NamePayOrder, constants.UsecaseLayer)
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
	if basketId == 0 {
		err = nil
	}
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	//если basketId-0, значит у пользователя нет корзины
	if basketId == 0 {
		//создаем заказ-корзину
		basketId, err = h.uc.Create(r.Context(), email)
		if err != nil {
			if err.Error() == repoErrors.CreateError {
				w = functions.ErrorResponse(w, repoErrors.CreateError, http.StatusInternalServerError)
				return
			}
			w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
			return
		}
	}
	//добавляем еду в заказ
	err = h.uc.AddFoodToOrder(r.Context(), alias.FoodId(foodId), basketId)
	if err != nil {
		if err.Error() == repoErrors.NotAddFood {
			w = functions.ErrorResponse(w, repoErrors.NotAddFood, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//PUT - ok

func (h *OrderHandler) UpdateFoodCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := r.Context().Value("email").(string)
	basketId, err := h.uc.GetBasketId(r.Context(), email)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if err = r.Body.Close(); err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	var item FoodCount
	err = json.Unmarshal(body, &item)
	fmt.Println("we are updating food count ", item.FoodId, item.Count)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	err = h.uc.UpdateCountInOrder(r.Context(), basketId, alias.FoodId(item.FoodId), uint32(item.Count))
	fmt.Println(err)
	if err != nil {
		if err.Error() == repoErrors.NotAddFood {
			w = functions.ErrorResponse(w, repoErrors.NotAddFood, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusUnauthorized)
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
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	err = h.uc.DeleteFromOrder(r.Context(), basketId, alias.FoodId(foodId))
	if err != nil {
		if err.Error() == repoErrors.NotDeleteFood {
			w = functions.ErrorResponse(w, repoErrors.NotDeleteFood, http.StatusInternalServerError)
			return
		}
		fmt.Println(err)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
