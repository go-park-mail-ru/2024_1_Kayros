package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"github.com/satori/uuid"
	"go.uber.org/zap"
)

type Delivery struct {
	ucSession session.Usecase
	ucUser    user.Usecase
	logger    *zap.Logger
}

func NewDeliveryLayer(ucSessionProps session.Usecase, ucUserProps user.Usecase, loggerProps *zap.Logger) *Delivery {
	return &Delivery{
		ucSession: ucSessionProps,
		ucUser:    ucUserProps,
		logger:    loggerProps,
	}
}

func (d *Delivery) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		err := errors.New("request_id передан не был")
		functions.LogError(d.logger, requestId, cnst.NameHandlerSignUp, err, cnst.DeliveryLayer)
	} else {
		requestId = ctxRequestId.(string)
	}
	email := r.Context().Value("email")
	if email != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignUp, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignUp, err, http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	var bodyDataDTO dto.SignUp
	err = json.Unmarshal(requestBody, &bodyDataDTO)
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignUp, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}

	isValid, err := bodyDataDTO.Validate()
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignUp, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	if !isValid {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignUp, errors.New(myerrors.BadCredentialsError), http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}

	isExist, err := d.ucUser.IsExistByEmail(r.Context(), bodyDataDTO.Email)
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignUp, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if isExist {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignUp, errors.New(myerrors.UserAlreadyExistError), http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UserAlreadyExistError, http.StatusBadRequest)
		return
	}

	u := dto.NewUserFromSignUp(&bodyDataDTO)
	u, err = d.ucUser.Create(r.Context(), u)
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignUp, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	sessionId := uuid.NewV4().String()
	err = d.ucSession.SetValue(r.Context(), alias.SessionKey(sessionId), alias.SessionValue(bodyDataDTO.Email))
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignUp, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	// собираем Cookie
	expiration := time.Now().Add(14 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Expires:  expiration,
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)

	uDTO := dto.NewUser(u)
	w = functions.JsonResponse(w, uDTO)
	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerSignUp, cnst.DeliveryLayer)
}

func (d *Delivery) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := r.Context().Value("request_id").(string)
	email := r.Context().Value("email")
	if email != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignIn, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignIn, err, http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	var bodyDTO dto.SignIn
	err = json.Unmarshal(body, &bodyDTO)
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignIn, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}

	isValid, err := bodyDTO.Validate()
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignIn, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	if !isValid {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignIn, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}

	u, err := d.ucUser.GetByEmail(r.Context(), bodyDTO.Email)
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignIn, err, http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if u == nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignIn, errors.New(myerrors.BadAuthCredentialsError), http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadAuthCredentialsError, http.StatusBadRequest)
		return
	}

	isEqual, err := d.ucUser.CheckPassword(r.Context(), bodyDTO.Email, bodyDTO.Password)
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignIn, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	if isEqual {
		sessionId := uuid.NewV4()
		expiration := time.Now().Add(14 * 24 * time.Hour)
		cookie := http.Cookie{
			Name:     "session_id",
			Value:    sessionId.String(),
			Expires:  expiration,
			HttpOnly: false,
		}
		http.SetCookie(w, &cookie)

		err = d.ucSession.SetValue(r.Context(), alias.SessionKey(sessionId.String()), alias.SessionValue(u.Email))
		if err != nil {
			functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignIn, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
			w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
			return
		}

		// Собираем ответ
		uDTO := dto.NewUser(u)
		w = functions.JsonResponse(w, uDTO)
		return
	}
	w = functions.ErrorResponse(w, myerrors.BadAuthCredentialsError, http.StatusBadRequest)
	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerSignIn, cnst.DeliveryLayer)
}

func (d *Delivery) SignOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := r.Context().Value("request_id").(string)
	email := r.Context().Value("email")
	if email == nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignOut, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignOut, err, http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	wasDeleted, err := d.ucSession.DeleteKey(r.Context(), alias.SessionKey(sessionCookie.Value))
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignOut, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if !wasDeleted {
		functions.LogWarn(d.logger, requestId, cnst.NameHandlerSignOut, errors.New("Такого ключа нет в Redis"), cnst.DeliveryLayer)
	}

	sessionCookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, sessionCookie)

	w = functions.JsonResponse(w, "Сессия успешно завершена")
	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerSignUp, cnst.DeliveryLayer)
}
