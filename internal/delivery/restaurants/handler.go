package delivery

import (
	"encoding/json"
	"net/http"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/entity/dto"
	rest "2024_1_kayros/internal/usecase/restaurants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
)

type RestaurantHandler struct {
	uc *rest.RestaurantUseCase
}

func NewRestaurantHandler(h *rest.RestaurantUseCase) *RestaurantHandler {
	return &RestaurantHandler{uc: h}
}

func (h *RestaurantHandler) RestaurantList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rests []*entity.Restaurant
	rests, err := h.uc.GetAll(r.Context())
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
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
