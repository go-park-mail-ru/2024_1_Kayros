package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
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
	fmt.Print("we are in restaurants")
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
	body, err := json.Marshal(restsDTO)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodGetAllRests, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodGetAllRests, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	functions.LogOk(h.logger, requestId, constants.NameMethodGetAllRests, constants.DeliveryLayer)
}

func (h *RestaurantHandler) RestaurantById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		err := errors.New("request_id передан не был")
		functions.LogError(h.logger, requestId, constants.NameMethodGetRestById, err, constants.DeliveryLayer)
	} else {
		requestId = ctxRequestId.(string)
	}
	rest, err := h.ucRest.GetById(r.Context(), alias.RestId(id))
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodGetRestById, err, constants.DeliveryLayer)
		if err.Error() == restaurants.NoRestError {
			w = functions.ErrorResponse(w, restaurants.NoRestError, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	categories, err := h.ucFood.GetByRestId(r.Context(), alias.RestId(id))
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodGetFoodByRest, err, constants.DeliveryLayer)
		if err.Error() == foodRepo.NoFoodError {
			fmt.Println()
			w = functions.ErrorResponse(w, foodRepo.NoFoodError, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	restDTO := dto.NewRestaurantAndFood(rest, categories)
	body, err := json.Marshal(restDTO)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodGetRestById, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		functions.LogError(h.logger, requestId, constants.NameMethodGetRestById, err, constants.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	functions.LogOk(h.logger, requestId, constants.NameMethodGetRestById, constants.DeliveryLayer)
}
