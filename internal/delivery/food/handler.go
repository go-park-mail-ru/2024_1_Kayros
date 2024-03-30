package delivery

import (
	"encoding/json"
	"io"
	"net/http"

	food "2024_1_kayros/internal/usecase/food"
	"2024_1_kayros/internal/utils/functions"
)

type FoodInOrder struct {
	OrderId int `json:"order_id"`
	FoodId  int `json:"food_id"`
	Count   int `json:"count"`
}

type FoodHandler struct {
	uc *food.UseCase
}

func NewFoodHandler(u *food.UseCase) *FoodHandler {
	return &FoodHandler{uc: u}
}

func (h *FoodHandler) UpdateCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	var fio FoodInOrder
	err = json.Unmarshal(body, &fio)
	if err != nil {
		w = functions.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = h.uc.UpdateCountInOrder(r.Context(), fio.OrderId, fio.FoodId, fio.Count)
	if err != nil {
		w = functions.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
