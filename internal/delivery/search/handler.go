package search

import (
	"net/http"

	"go.uber.org/zap"

	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/search"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
)

type Delivery struct {
	ucSearch search.Usecase
	logger   *zap.Logger
}

func NewDelivery(ucs search.Usecase, loggerProps *zap.Logger) *Delivery {
	return &Delivery{
		ucSearch: ucs,
		logger:   loggerProps,
	}
}

func (h *Delivery) Search(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	str := r.URL.Query().Get("search")
	var rests []*dto.RestaurantAndFood
	var err error
	if str != "" {
		rests, err = h.ucSearch.Search(r.Context(), str)
		if err != nil {
			h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
	} else {
		return
	}
	restaurantAndFoodArray := &dto.RestaurantAndFoodArray{Payload: rests}
	w = functions.JsonResponse(w, restaurantAndFoodArray)
}
