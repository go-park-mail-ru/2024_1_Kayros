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

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/entity/dto"
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
		h.logger.Error("request_id передан не был")
	} else {
		requestId = ctxRequestId.(string)
	}

	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	token := ""
	ctxToken := r.Context().Value("unauth_token")
	if ctxToken != nil {
		token = ctxToken.(string)
	}
	if email == "" && token == "" {
		h.logger.Error(myerrors.NoBasket.Error())
		return
	}

	// неавторизованный - вовзвращаем по токену, авторизованный - по почте
	//неавторизованный пользователь
	var order *entity.Order
	var err error
	if token != "" {
		fmt.Println(email)
		order, err = h.uc.GetBasketNoAuth(r.Context(), token)
	} else if email != "" {
		fmt.Println(token)
		order, err = h.uc.GetBasket(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoBasket, http.StatusOK)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	orderDTO := dto.NewOrder(order)
	w = functions.JsonResponse(w, orderDTO)
}

//PUT

func (h *OrderHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		h.logger.Error("request_id передан не был")
	} else {
		requestId = ctxRequestId.(string)
	}

	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	token := ""
	ctxToken := r.Context().Value("unauth_token")
	if ctxToken != nil {
		token = ctxToken.(string)
	}
	if email == "" && token == "" {
		h.logger.Error(myerrors.NoBasket.Error())
		return
	}

	var fullAddress dto.FullAddress
	var basketId alias.OrderId
	var err error
	if token != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), token)
	} else if email != "" {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil {
		//h.logger.Error(myerrors.NoBasket.Error(), zap.String(constants.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			//h.logger.Error(myerrors.NoBasket.Error(), zap.String(constants.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.NoBasket, http.StatusOK)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	if err = r.Body.Close(); err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &fullAddress)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	if len(fullAddress.Address) < 14 || len(fullAddress.ExtraAddress) < 2 {
		h.logger.Error(myerrors.InvalidAddressEn.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InvalidAddress, http.StatusBadRequest)
		return
	}

	err = h.uc.UpdateAddress(r.Context(), fullAddress, basketId)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//PUT - ok

func (h *OrderHandler) Pay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		h.logger.Error("request_id передан не был")
	} else {
		requestId = ctxRequestId.(string)
	}

	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	if email == "" {
		h.logger.Error(myerrors.CtxEmail.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.NoAuth, http.StatusUnauthorized)
		return
	}
	token := ""
	ctxToken := r.Context().Value("unauth_token")
	if ctxToken != nil {
		token = ctxToken.(string)
	}
	var basket *entity.Order
	var err error
	if token != "" {
		// эта корзина создана неавторизованным
		basket, err = h.uc.GetBasketNoAuth(r.Context(), token)
	} else {
		// в случае, когда заказ был создан авторизированным пользователем
		basket, err = h.uc.GetBasket(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoBasket, http.StatusOK)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	payedOrder, err := h.uc.Pay(r.Context(), alias.OrderId(basket.Id), basket.Status, email, alias.UserId(basket.UserId))
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		if errors.Is(err, myerrors.AlreadyPayed) {
			w = functions.ErrorResponse(w, myerrors.AlreadyPayed, http.StatusBadRequest)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusUnauthorized)
		return
	}

	response := payedOrderInfo{Id: alias.OrderId(payedOrder.Id), Status: payedOrder.Status}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// POST- ok

func (h *OrderHandler) AddFood(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		h.logger.Error("request_id передан не был")
	} else {
		requestId = ctxRequestId.(string)
	}

	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	token := ""
	ctxToken := r.Context().Value("unauth_token")
	if ctxToken != nil {
		token = ctxToken.(string)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	if err = r.Body.Close(); err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	var item FoodCount
	err = json.Unmarshal(body, &item)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	var basketId alias.OrderId
	if token != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), token)

	} else if email != "" {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil && !errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	//если basketId-0, значит у пользователя нет корзины
	if basketId == 0 {
		//создаем заказ-корзину
		if token != "" {
			basketId, err = h.uc.CreateNoAuth(r.Context(), token)
		} else if email != "" {
			basketId, err = h.uc.Create(r.Context(), email)
		}
		if err != nil {
			h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
	}

	//добавляем еду в заказ
	err = h.uc.AddFoodToOrder(r.Context(), item.FoodId, item.Count, basketId)
	if err != nil {
		if errors.Is(err, myerrors.OrderAddFood) {
			h.logger.Error(myerrors.OrderAddFood.Error(), zap.String(constants.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.NoAddFoodToOrder, http.StatusInternalServerError)
			return
		}
		h.logger.Error(myerrors.OrderAddFood.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	order, err := h.uc.GetOrderById(r.Context(), basketId)
	if err != nil {
		h.logger.Error(myerrors.OrderAddFood.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	orderDTO := dto.NewOrder(order)
	w = functions.JsonResponse(w, orderDTO)
}

func (h *OrderHandler) Clean(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		h.logger.Error("request_id передан не был")
	} else {
		requestId = ctxRequestId.(string)
	}

	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	token := ""
	ctxToken := r.Context().Value("unauth_token")
	if ctxToken != nil {
		token = ctxToken.(string)
	}
	if email == "" && token == "" {
		w = functions.ErrorResponse(w, myerrors.NoBasket, http.StatusBadRequest)
		return
	}

	var basketId alias.OrderId
	var err error
	if token != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), token)
	} else if email != "" {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoBasket, http.StatusOK)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	err = h.uc.Clean(r.Context(), basketId)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		if errors.Is(err, myerrors.FailCleanBasket) {
			w = functions.ErrorResponse(w, myerrors.NoClean, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	w = functions.ErrorResponse(w, myerrors.SuccessClean, http.StatusOK)
}

//PUT - ok

func (h *OrderHandler) UpdateFoodCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		h.logger.Error("request_id передан не был")
	} else {
		requestId = ctxRequestId.(string)
	}

	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	token := ""
	ctxToken := r.Context().Value("unauth_token")
	if ctxToken != nil {
		token = ctxToken.(string)
	}

	var basketId alias.OrderId
	var err error
	if token != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), token)
	} else if email != "" {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoBasket, http.StatusOK)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	if err = r.Body.Close(); err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	var item FoodCount
	err = json.Unmarshal(body, &item)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
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
			w = functions.ErrorResponse(w, myerrors.NoAddFoodToOrder, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	orderDTO := dto.NewOrder(order)
	w = functions.JsonResponse(w, orderDTO)
	w.WriteHeader(http.StatusOK)
}

//DELETE

func (h *OrderHandler) DeleteFoodFromOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	foodId, _ := strconv.Atoi(vars["food_id"])
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		h.logger.Error("request_id передан не был")
	} else {
		requestId = ctxRequestId.(string)
	}

	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	token := ""
	ctxToken := r.Context().Value("unauth_token")
	if ctxToken != nil {
		token = ctxToken.(string)
	}

	var basketId alias.OrderId
	var err error
	if token != "" {
		basketId, err = h.uc.GetBasketIdNoAuth(r.Context(), token)
	} else if email != "" {
		basketId, err = h.uc.GetBasketId(r.Context(), email)
	}
	if err != nil {
		h.logger.Error(err.Error(), zap.String(constants.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoBasket, http.StatusOK)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	order, err := h.uc.DeleteFromOrder(r.Context(), basketId, alias.FoodId(foodId))
	if err != nil {
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoDeleteFood, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	orderDTO := dto.NewOrder(order)
	w = functions.JsonResponse(w, orderDTO)
	w.WriteHeader(http.StatusOK)
}
