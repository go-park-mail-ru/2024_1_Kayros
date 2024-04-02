package restaurants

import (
	"net/http"

	"2024_1_kayros/internal/entity/dto"
	rest "2024_1_kayros/internal/usecase/restaurants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
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
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
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
