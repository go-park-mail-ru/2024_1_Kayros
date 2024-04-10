package delivery

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"2024_1_kayros/internal/entity/dto"
	foodRepo "2024_1_kayros/internal/repository/food"
	"2024_1_kayros/internal/repository/restaurants"
	foodUc "2024_1_kayros/internal/usecase/food"
	restUc "2024_1_kayros/internal/usecase/restaurants"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/constants"
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
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		err := errors.New("request_id передан не был")
		functions.LogError(h.logger, requestId, constants.NameMethodGetAllRests, err, constants.DeliveryLayer)
	} else {
		requestId = ctxRequestId.(string)
	}
	rests, err := h.ucRest.GetAll(r.Context())
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodGetAllRests, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	restsDTO := dto.NewRestaurantArray(rests)
	w = functions.JsonResponse(w, restsDTO)
	functions.LogOk(h.logger, requestId, constants.NameMethodGetAllRests, constants.DeliveryLayer)
}

func (h *RestaurantHandler) RestaurantById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		err := errors.New("request_id передан не был")
		functions.LogError(h.logger, requestId, constants.NameMethodGetRestById, err, constants.DeliveryLayer)
	} else {
		requestId = ctxRequestId.(string)
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodGetRestById, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.NotFoundError, http.StatusNotFound)
		return
	}
	rest, err := h.ucRest.GetById(r.Context(), alias.RestId(id))
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodGetRestById, err, constants.DeliveryLayer)
		if err.Error() == restaurants.NoRestError {
			w = functions.ErrorResponse(w, restaurants.NoRestError, http.StatusNotFound)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	categories, err := h.ucFood.GetByRestId(r.Context(), alias.RestId(id))
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodGetFoodByRest, err, constants.DeliveryLayer)
		if err.Error() == foodRepo.NoFoodError {
			w = functions.ErrorResponse(w, foodRepo.NoFoodError, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	restDTO := dto.NewRestaurantAndFood(rest, categories)
	w = functions.JsonResponse(w, restDTO)
	functions.LogOk(h.logger, requestId, constants.NameMethodGetRestById, constants.DeliveryLayer)
}
