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
	FoodId alias.FoodId `json:"food_id"`
	Count  uint32       `json:"count"`
}

type payedOrderInfo struct {
	Id     alias.OrderId `json:"id"`
	Status string        `json:"status"`
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
	}
	if email == "" {
		functions.LogErrorResponse(h.logger, requestId, constants.NameMethodGetBasket, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}
	order, err := h.uc.GetBasket(r.Context(), email)
	if err != nil {
		fmt.Println(err)
		if err.Error() == "У Вас нет корзины" {
			functions.LogInfo(h.logger, requestId, constants.NameMethodGetBasket, repoErrors.NoBasketError, constants.DeliveryLayer)
			w = functions.ErrorResponse(w, repoErrors.NoBasketError, http.StatusOK)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	orderDTO := dto.NewOrder(order)
	w = functions.JsonResponse(w, orderDTO)
	functions.LogOk(h.logger, requestId, constants.NameMethodGetBasket, constants.DeliveryLayer)
}

//PUT - ok

func (h *OrderHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
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
	}
	if email == "" {
		functions.LogErrorResponse(h.logger, requestId, constants.NameMethodUpdateAddress, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}
	var fullAddress dto.FullAddress
	basketId, err := h.uc.GetBasketId(r.Context(), email)
	if err != nil {
		functions.LogWarn(h.logger, requestId, constants.NameMethodGetFoodById, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if basketId == 0 {
		w = functions.ErrorResponse(w, repoErrors.NoBasketError, http.StatusOK)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		functions.LogErrorResponse(h.logger, requestId, constants.NameMethodUpdateAddress, err, http.StatusInternalServerError, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if err = r.Body.Close(); err != nil {
		functions.LogErrorResponse(h.logger, requestId, constants.NameMethodUpdateAddress, err, http.StatusInternalServerError, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &fullAddress)
	if err != nil {
		functions.LogErrorResponse(h.logger, requestId, constants.NameMethodUpdateAddress, err, http.StatusInternalServerError, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	if len(fullAddress.Address) < 14 || len(fullAddress.ExtraAddress) < 2 {
		w = functions.ErrorResponse(w, "Некорректный адрес", http.StatusBadRequest)
		return
	}
	err = h.uc.UpdateAddress(r.Context(), fullAddress, basketId)
	if err != nil {
		if errors.Is(err, fmt.Errorf(repoErrors.NotUpdateError)) {
			functions.LogErrorResponse(h.logger, requestId, constants.NameMethodUpdateAddress, fmt.Errorf(repoErrors.NotUpdateError), http.StatusInternalServerError, constants.DeliveryLayer)
			w = functions.ErrorResponse(w, repoErrors.NotUpdateError, http.StatusInternalServerError)
			return
		}
		functions.LogErrorResponse(h.logger, requestId, constants.NameMethodUpdateAddress, err, http.StatusInternalServerError, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	functions.LogOk(h.logger, requestId, constants.NameMethodUpdateAddress, constants.DeliveryLayer)
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
	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	if email == "" {
		functions.LogErrorResponse(h.logger, requestId, constants.NameMethodPayOrder, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}
	basket, err := h.uc.GetBasket(r.Context(), email)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodGetBasket, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if len(basket.Food) == 0 || basket.Id == 0 {
		functions.LogWarn(h.logger, requestId, constants.NameMethodGetBasket, fmt.Errorf(repoErrors.EmptyError), constants.DeliveryLayer)
		w = functions.ErrorResponse(w, repoErrors.EmptyError, http.StatusOK)
		return
	}
	payedOrder, err := h.uc.Pay(r.Context(), alias.OrderId(basket.Id), basket.Status)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodPayOrder, err, constants.DeliveryLayer)
		if errors.Is(err, fmt.Errorf(repoErrors.NotUpdateStatusError)) {
			w = functions.ErrorResponse(w, repoErrors.NotUpdateStatusError, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusUnauthorized)
		return
	}
	response := payedOrderInfo{Id: alias.OrderId(payedOrder.Id), Status: payedOrder.Status}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodPayOrder, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodPayOrder, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	functions.LogOk(h.logger, requestId, constants.NameMethodPayOrder, constants.DeliveryLayer)
}

// POST-ok

func (h *OrderHandler) AddFood(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		err := errors.New("request_id передан не был")
		functions.LogError(h.logger, requestId, constants.NameMethodAddFood, err, constants.DeliveryLayer)
	} else {
		requestId = ctxRequestId.(string)
	}
	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	if email == "" {
		functions.LogErrorResponse(h.logger, requestId, constants.NameMethodAddFood, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodUpdateCountInOrder, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if err = r.Body.Close(); err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodUpdateCountInOrder, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	var item FoodCount
	err = json.Unmarshal(body, &item)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodUpdateCountInOrder, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	basketId, err := h.uc.GetBasketId(r.Context(), email)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodAddFood, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	//если basketId-0, значит у пользователя нет корзины
	if basketId == 0 {
		//создаем заказ-корзину
		basketId, err = h.uc.Create(r.Context(), email)
		if err != nil {
			functions.LogError(h.logger, requestId, constants.NameMethodAddFood, err, constants.DeliveryLayer)
			if errors.Is(err, fmt.Errorf(repoErrors.CreateError)) {
				w = functions.ErrorResponse(w, repoErrors.CreateError, http.StatusInternalServerError)
			} else {
				w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
			}
			return
		}
	}
	//добавляем еду в заказ
	err = h.uc.AddFoodToOrder(r.Context(), item.FoodId, item.Count, basketId)
	if err != nil {
		if errors.Is(err, fmt.Errorf(repoErrors.NotAddFood)) {
			functions.LogError(h.logger, requestId, constants.NameMethodAddFood, fmt.Errorf(repoErrors.NotAddFood), constants.DeliveryLayer)
			w = functions.ErrorResponse(w, repoErrors.NotAddFood, http.StatusInternalServerError)
			return
		}
		functions.LogError(h.logger, requestId, constants.NameMethodAddFood, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	order, err := h.uc.GetBasket(r.Context(), email)
	if err != nil {
		functions.LogErrorResponse(h.logger, requestId, constants.NameMethodGetBasket, fmt.Errorf(repoErrors.NoBasketError), http.StatusNotFound, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if len(order.Food) != 0 {
		order.RestaurantId = order.Food[0].RestaurantId
	}
	orderDTO := dto.NewOrder(order)
	w = functions.JsonResponse(w, orderDTO)
	functions.LogOk(h.logger, requestId, constants.NameMethodAddFood, constants.DeliveryLayer)
}

func (h *OrderHandler) Clean(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	if email == "" {
		fmt.Println("Необходимо авторизиризоваться")
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}
	basketId, err := h.uc.GetBasketId(r.Context(), email)
	if err != nil {
		fmt.Println(err)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	//если basketId-0, значит у пользователя нет корзины
	if basketId == 0 {
		fmt.Println("Нет корзины, добавьте что-нибудь")
		w = functions.ErrorResponse(w, repoErrors.NoBasketError, http.StatusInternalServerError)
		return
	}
	err = h.uc.Clean(r.Context(), basketId)
	if err != nil {
		fmt.Println(err)
		w = functions.ErrorResponse(w, repoErrors.CleanError, http.StatusInternalServerError)
		return
	}
	w = functions.ErrorResponse(w, "Корзина очищена", http.StatusOK)
}

//PUT - ok

func (h *OrderHandler) UpdateFoodCount(w http.ResponseWriter, r *http.Request) {
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
	}
	if email == "" {
		functions.LogErrorResponse(h.logger, requestId, constants.NameMethodUpdateCountInOrder, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}
	basketId, err := h.uc.GetBasketId(r.Context(), email)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodUpdateCountInOrder, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if basketId == 0 {
		functions.LogWarn(h.logger, requestId, constants.NameMethodGetBasket, fmt.Errorf(repoErrors.EmptyError), constants.DeliveryLayer)
		w = functions.ErrorResponse(w, repoErrors.EmptyError, http.StatusOK)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodUpdateCountInOrder, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if err = r.Body.Close(); err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodUpdateCountInOrder, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	var item FoodCount
	err = json.Unmarshal(body, &item)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodUpdateCountInOrder, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	if item.FoodId <= 0 || item.Count <= 0 {
		functions.LogError(h.logger, requestId, constants.NameMethodUpdateCountInOrder, fmt.Errorf("id или кол-во ды отрицательное"), constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	order, err := h.uc.UpdateCountInOrder(r.Context(), basketId, alias.FoodId(item.FoodId), uint32(item.Count))
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodUpdateCountInOrder, err, constants.DeliveryLayer)
		if errors.Is(err, fmt.Errorf(repoErrors.NotAddFood)) {
			w = functions.ErrorResponse(w, repoErrors.NotAddFood, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	orderDTO := dto.NewOrder(order)
	w = functions.JsonResponse(w, orderDTO)
	functions.LogOk(h.logger, requestId, constants.NameMethodUpdateCountInOrder, constants.DeliveryLayer)
	w.WriteHeader(http.StatusOK)
}

//DELETE - ok

func (h *OrderHandler) DeleteFoodFromOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	foodId, _ := strconv.Atoi(vars["food_id"])
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
	}
	if email == "" {
		functions.LogErrorResponse(h.logger, requestId, constants.NameMethodDeleteFromOrder, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}
	basketId, err := h.uc.GetBasketId(r.Context(), email)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	order, err := h.uc.DeleteFromOrder(r.Context(), basketId, alias.FoodId(foodId))
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodDeleteFromOrder, err, constants.DeliveryLayer)
		if errors.Is(err, fmt.Errorf(repoErrors.NotDeleteFood)) {
			w = functions.ErrorResponse(w, repoErrors.NotDeleteFood, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	orderDTO := dto.NewOrder(order)
	w = functions.JsonResponse(w, orderDTO)
	w.WriteHeader(http.StatusOK)
	functions.LogOk(h.logger, requestId, constants.NameMethodDeleteFromOrder, constants.DeliveryLayer)
}
