package restaurants

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/entity/dto"
	foodUc "2024_1_kayros/internal/usecase/food"
	restUc "2024_1_kayros/internal/usecase/restaurants"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
)

type RestaurantAndFoodDTO struct {
	Id              uint64      `json:"id" valid:"-"`
	Name            string      `json:"name" valid:"-"`
	LongDescription string      `json:"long_description" valid:"-"`
	ImgUrl          string      `json:"img_url" valid:"url"`
	Food            []*dto.Food `json:"food"`
}

type RestaurantHandler struct {
	ucRest restUc.Usecase
	ucFood foodUc.Usecase
	logger *zap.Logger
}

func NewRestaurantHandler(ucr restUc.Usecase, ucf foodUc.Usecase, loggerProps *zap.Logger) *RestaurantHandler {
	return &RestaurantHandler{
		ucRest: ucr,
		ucFood: ucf,
		logger: loggerProps,
	}
}

func (h *RestaurantHandler) RestaurantList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	filter := r.URL.Query().Get("filter")
	var id = 0
	if filter != "" {
		id, err := strconv.Atoi(filter)
		if err != nil || id < 0 {
			h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
			return
		}
	}
	var rests []*entity.Restaurant
	var err error
	if filter != "" {
		rests, err = h.ucRest.GetByFilter(r.Context(), alias.CategoryId(id))
		if err != nil {
			h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
	} else {
		rests, err = h.ucRest.GetAll(r.Context())
		if err != nil {
			h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
	}

	restsDTO := dto.NewRestaurantArray(rests)
	w = functions.JsonResponse(w, restsDTO)
}

func (h *RestaurantHandler) RestaurantById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.NotFoundRu, http.StatusNotFound)
		return
	}
	rest, err := h.ucRest.GetById(r.Context(), alias.RestId(id))
	if err != nil {
		if errors.Is(err, myerrors.SqlNoRowsRestaurantRelation) {
			h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.NotFoundRu, http.StatusNotFound)
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	categories, err := h.ucFood.GetByRestId(r.Context(), alias.RestId(id))
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	restDTO := dto.NewRestaurantAndFood(rest, categories)
	w = functions.JsonResponse(w, restDTO)
}

func (h *RestaurantHandler) CategoryList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	categories, err := h.ucRest.GetCategoryList(r.Context())
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	catsDTO := dto.NewCategoryArray(categories)
	w = functions.JsonResponse(w, catsDTO)
}
