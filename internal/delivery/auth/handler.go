package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/auth"
	"2024_1_kayros/internal/usecase/session"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/sanitizer"

	"go.uber.org/zap"
)

type Delivery struct {
	ucSession session.Usecase // methods for communicating to microservice session
	ucAuth    auth.Usecase
	logger    *zap.Logger
	cfg       *config.Project
	metrics   *metrics.Metrics
}

func NewDeliveryLayer(cfgProps *config.Project, ucSessionProps session.Usecase, ucAuthProps auth.Usecase, loggerProps *zap.Logger, metrics *metrics.Metrics) *Delivery {
	return &Delivery{
		ucSession: ucSessionProps,
		ucAuth:    ucAuthProps,
		logger:    loggerProps,
		cfg:       cfgProps,
		metrics:   metrics,
	}
}

func (d *Delivery) SignUp(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email != "" {
		d.logger.Error(myerrors.CtxEmail.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.RegisteredRu, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	var signupDTO dto.UserSignUp
	err = json.Unmarshal(body, &signupDTO)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	isValid, err := signupDTO.Validate()
	if err != nil || !isValid {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	u := dto.NewUserFromSignUpForm(&signupDTO)
	uSignedUp, err := d.ucAuth.SignUp(r.Context(), u)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.UserAlreadyExist) {
			w = functions.ErrorResponse(w, myerrors.UserAlreadyExistRu, http.StatusBadRequest)
			return
		}
		// error `myerrors.SqlNoRowsUserRelation` is handled
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	uDTO := dto.NewUserData(sanitizer.User(uSignedUp))

	w, err = functions.SetCookie(w, r, d.ucSession, u.Email, d.cfg)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	w = functions.JsonResponse(w, uDTO)
}

func (d *Delivery) SignIn(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email != "" {
		d.logger.Error(myerrors.CtxEmail.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.AuthorizedRu, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	var bodyDTO dto.UserSignIn
	err = json.Unmarshal(body, &bodyDTO)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	isValid, err := bodyDTO.Validate()
	if err != nil || !isValid {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	u, err := d.ucAuth.SignIn(r.Context(), bodyDTO.Email, bodyDTO.Password)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) || errors.Is(err, myerrors.BadAuthPassword) {
			w = functions.ErrorResponse(w, myerrors.BadAuthCredentialsRu, http.StatusBadRequest)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	uDTO := dto.NewUserData(sanitizer.User(u))

	w, err = functions.SetCookie(w, r, d.ucSession, u.Email, d.cfg)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	w = functions.JsonResponse(w, uDTO)
}

func (d *Delivery) SignOut(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		w = functions.ErrorResponse(w, myerrors.SignOutAlreadyRu, http.StatusUnauthorized)
		return
	}

	w, err := functions.FlashCookie(r, w, d.ucSession, &d.cfg.Redis, d.metrics)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	w = functions.JsonResponse(w, &dto.ResponseDetail{Detail: "Сессия успешно завершена"})
}
