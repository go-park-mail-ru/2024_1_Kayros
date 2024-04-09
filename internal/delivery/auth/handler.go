package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"2024_1_kayros/config"
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
	ucCsrf    session.Usecase
	logger    *zap.Logger
	cfg       *config.Project
}

func NewDeliveryLayer(cfgProps *config.Project, ucSessionProps session.Usecase, ucUserProps user.Usecase, ucCsrfProps session.Usecase, loggerProps *zap.Logger) *Delivery {
	return &Delivery{
		ucSession: ucSessionProps,
		ucUser:    ucUserProps,
		ucCsrf:    ucCsrfProps,
		logger:    loggerProps,
		cfg:       cfgProps,
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

	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	if email != "" {
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
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignUp, errors.New(myerrors.UserAlreadyExistError), http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UserAlreadyExistError, http.StatusUnauthorized)
		return
	}
	if isExist {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignUp, errors.New(myerrors.UserAlreadyExistError), http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UserAlreadyExistError, http.StatusUnauthorized)
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
	expiration := time.Now().Add(14 * 24 * time.Hour)
	sessionCookie := http.Cookie{
		Name:     cnst.SessionCookieName,
		Value:    sessionId,
		Expires:  expiration,
		HttpOnly: false,
	}

	csrfToken, err := genCsrfToken(d.logger, requestId, cnst.NameHandlerSignUp, d.cfg.CsrfSecretKey, alias.SessionKey(sessionId))
	if err == nil {
		err = d.ucCsrf.SetValue(r.Context(), alias.SessionKey(csrfToken), alias.SessionValue(u.Email))
		if err != nil {
			functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignUp, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
			w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
			return
		}
		csrfCookie := http.Cookie{
			Name:     cnst.CsrfCookieName,
			Value:    csrfToken,
			Expires:  expiration,
			HttpOnly: false,
		}
		http.SetCookie(w, &csrfCookie)
	}
	http.SetCookie(w, &sessionCookie)
	uDTO := dto.NewUser(u)
	w = functions.JsonResponse(w, uDTO)
	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerSignUp, cnst.DeliveryLayer)
}

func (d *Delivery) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		err := errors.New("request_id передан не был")
		functions.LogError(d.logger, requestId, cnst.NameHandlerSignUp, err, cnst.DeliveryLayer)
	} else {
		requestId = ctxRequestId.(string)
	}

	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	if email != "" {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignUp, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, cnst.DeliveryLayer)
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
			Name:     cnst.SessionCookieName,
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

		csrfToken, err := genCsrfToken(d.logger, requestId, cnst.NameHandlerSignUp, d.cfg.CsrfSecretKey, alias.SessionKey(sessionId.String()))
		if err == nil {
			err = d.ucCsrf.SetValue(r.Context(), alias.SessionKey(csrfToken), alias.SessionValue(u.Email))
			if err != nil {
				functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignUp, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
				w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
				return
			}
			csrfCookie := http.Cookie{
				Name:     cnst.CsrfCookieName,
				Value:    csrfToken,
				Expires:  expiration,
				HttpOnly: false,
			}
			http.SetCookie(w, &csrfCookie)
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
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		err := errors.New("request_id передан не был")
		functions.LogError(d.logger, requestId, cnst.NameHandlerSignUp, err, cnst.DeliveryLayer)
	} else {
		requestId = ctxRequestId.(string)
	}

	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	if email == "" {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignUp, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, cnst.DeliveryLayer)
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

	//

	csrfCookie, err := r.Cookie("csrf_token")
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignOut, err, http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	wasDeleted, err = d.ucSession.DeleteKey(r.Context(), alias.SessionKey(csrfCookie.Value))
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignOut, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if !wasDeleted {
		functions.LogWarn(d.logger, requestId, cnst.NameHandlerSignOut, errors.New("Такого ключа нет в Redis"), cnst.DeliveryLayer)
	}

	w = functions.JsonResponse(w, map[string]string{"detail": "Сессия успешно завершена"})
	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerSignUp, cnst.DeliveryLayer)
}

func genCsrfToken(logger *zap.Logger, requestId string, methodName string, secretKey string, sessionId alias.SessionKey) (string, error) {
	// Создание csrf_token
	hashData, err := functions.HashCsrf(secretKey, string(sessionId))
	if err != nil {
		functions.LogError(logger, requestId, methodName, err, cnst.DeliveryLayer)
		return "", err
	}
	csrfToken := hashData + "." + string(sessionId)
	return csrfToken, nil
}
