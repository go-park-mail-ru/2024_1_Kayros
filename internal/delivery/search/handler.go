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

type SearchDelivery struct {
	ucSearch search.Usecase
	logger   *zap.Logger
}

func NewSearchDelivery(ucs search.Usecase, loggerProps *zap.Logger) *SearchDelivery {
	return &SearchDelivery{
		ucSearch: ucs,
		logger:   loggerProps,
	}
}

func (h *SearchDelivery) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	w = functions.JsonResponse(w, rests)
}
