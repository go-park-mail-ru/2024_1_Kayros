package restaurants

import (
	"encoding/json"
	"net/http"
	"strconv"

	rest "2024_1_kayros/internal/usecase/restaurants"
	"2024_1_kayros/internal/utils/myerrors"
	"github.com/gorilla/mux"

	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/utils/functions"
)

type Delivery struct {
	ucRest rest.Usecase
}

func NewDelivery(ucRestProps rest.Usecase) *Delivery {
	return &Delivery{ucRest: ucRestProps}
}

func (h *Delivery) RestaurantList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rests, err := h.ucRest.GetAll(r.Context())
	if err != nil {
		w = functions.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}
	restsDTO := make([]*dto.RestaurantDTO, 0, len(rests)+1)
	for i, r := range rests {
		restsDTO[i].Id = r.Id
		restsDTO[i].Name = r.Name
		restsDTO[i].ShortDescription = r.ShortDescription
		restsDTO[i].ImgUrl = r.ImgUrl
	}
	w = functions.JsonResponse(w, restsDTO)
}

func (h *Delivery) RestaurantById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		functions.ErrorResponse(w, myerrors.NotFoundError)
	}
	rest, err := h.ucRest.GetById(r.Context(), id)
	if err != nil {
		w = functions.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}
	food, err := h.ucFood.GetByRest(r.Context(), uint64(id))
	var foodDTO []*dto.FoodDTO
	for i := range food {
		foodDTO[i].Id = food[i].Id
		foodDTO[i].Name = food[i].Name
		foodDTO[i].Description = food[i].Description
		foodDTO[i].ImgUrl = food[i].ImgUrl
		foodDTO[i].Weight = food[i].Weight
		foodDTO[i].Price = food[i].Price
		foodDTO[i].Restaurant = food[i].Restaurant
	}
	restDTO := &RestaurantAndFoodDTO{
		Id:              rest.Id,
		Name:            rest.Name,
		LongDescription: rest.LongDescription,
		ImgUrl:          rest.ImgUrl,
		Food:            foodDTO,
	}
	body, err := json.Marshal(restDTO)
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
