package restaurants

import (
	"net/http"
	"strconv"

	"2024_1_kayros/internal/usecase/food"
	rest "2024_1_kayros/internal/usecase/restaurants"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/myerrors"
	"github.com/gorilla/mux"

	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/utils/functions"
)

type Delivery struct {
	ucRest rest.Usecase
	ucFood food.Usecase
}

func NewDelivery(ucRestProps rest.Usecase, ucFoodProps food.Usecase) *Delivery {
	return &Delivery{
		ucRest: ucRestProps,
		ucFood: ucFoodProps,
	}
}

func (d *Delivery) RestaurantList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rests, err := d.ucRest.GetAll(r.Context())
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	restsDTO := dto.NewRestaurantArray(rests)
	w = functions.JsonResponse(w, restsDTO)
}

func (d *Delivery) RestaurantById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	restId, err := strconv.Atoi(vars["id"])
	if err != nil {
		functions.ErrorResponse(w, myerrors.NotFoundError, http.StatusNotFound)
		return
	}

	restaurant, err := d.ucRest.GetById(r.Context(), alias.RestId(restId))
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	foodArray, err := d.ucFood.GetByRestId(r.Context(), alias.RestId(restaurant.Id))
	restDTO := dto.NewRestaurantAndFood(restaurant, foodArray)
	w = functions.JsonResponse(w, restDTO)
}
