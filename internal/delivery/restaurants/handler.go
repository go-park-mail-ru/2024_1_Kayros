package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/entity/dto"
	foodUc "2024_1_kayros/internal/usecase/food"
	restUc "2024_1_kayros/internal/usecase/restaurants"
	"2024_1_kayros/internal/utils/functions"
)

type RestaurantAndFoodDTO struct {
	Id              uint64         `json:"id" valid:"-"`
	Name            string         `json:"name" valid:"-"`
	LongDescription string         `json:"long_description" valid:"-"`
	ImgUrl          string         `json:"img_url" valid:"url"`
	Food            []*dto.FoodDTO `json:"food"`
}

type RestaurantHandler struct {
	ucRest *restUc.RestaurantUseCase
	ucFood *foodUc.UseCase
}

func NewRestaurantHandler(ucr *restUc.RestaurantUseCase, ucf *foodUc.UseCase) *RestaurantHandler {
	return &RestaurantHandler{ucRest: ucr, ucFood: ucf}
}

func (h *RestaurantHandler) RestaurantList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rests []*entity.Restaurant
	rests, err := h.ucRest.GetAll(r.Context())
	if err != nil {
		w = functions.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var restsDTO []*dto.RestaurantDTO
	for i := range rests {
		restsDTO[i].Id = rests[i].Id
		restsDTO[i].Name = rests[i].Name
		restsDTO[i].ShortDescription = rests[i].ShortDescription
		restsDTO[i].ImgUrl = rests[i].ImgUrl
	}
	body, err := json.Marshal(restsDTO)
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

func (h *RestaurantHandler) RestaurantById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	rest, err := h.ucRest.GetById(r.Context(), id)
	if err != nil {
		w = functions.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}
	food, err := h.ucFood.GetByRest(r.Context(), id)
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
