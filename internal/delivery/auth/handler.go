package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/satori/uuid"
	"go.uber.org/zap"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/auth"
	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/props"
	"2024_1_kayros/internal/utils/sanitizer"
<<<<<<< HEAD
=======
	"go.uber.org/zap"
>>>>>>> fix_csrf_test
)

type Delivery struct {
	ucSession session.Usecase
	ucUser    user.Usecase
	ucCsrf    session.Usecase
	ucAuth    auth.Usecase
	logger    *zap.Logger
	cfg       *config.Project
}

func NewDeliveryLayer(cfgProps *config.Project, ucSessionProps session.Usecase, ucUserProps user.Usecase, ucCsrfProps session.Usecase, ucAuthProps auth.Usecase, loggerProps *zap.Logger) *Delivery {
	return &Delivery{
		ucSession: ucSessionProps,
		ucUser:    ucUserProps,
		ucCsrf:    ucCsrfProps,
		ucAuth:    ucAuthProps,
		logger:    loggerProps,
		cfg:       cfgProps,
	}
}

func (d *Delivery) SignUp(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email != "" {
<<<<<<< HEAD
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
	http.SetCookie(w, &sessionCookie)

	csrfToken, err := GenCsrfToken(d.logger, requestId, cnst.NameHandlerSignUp, d.cfg.CsrfSecretKey, alias.SessionKey(sessionId))
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
	u = sanitizer.User(u)
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
=======
		w = functions.ErrorResponse(w, myerrors.RegisteredRu, http.StatusUnauthorized)
>>>>>>> fix_csrf_test
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

	uAuth, err := d.ucAuth.SignUpUser(r.Context(), u.Email, u)
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
	u = sanitizer.User(uAuth)
	uDTO := dto.NewUserData(u)

	setCookieProps := props.GetSetCookieProps(d.ucCsrf, d.ucSession, u.Email, d.cfg.CsrfSecretKey)
	w, err = functions.SetCookie(w, r, setCookieProps)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	w = functions.JsonResponse(w, uDTO)
}

func (d *Delivery) SignIn(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email != "" {
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

	u, err := d.ucAuth.SignInUser(r.Context(), bodyDTO.Email, bodyDTO.Password)
	if err != nil {
<<<<<<< HEAD
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
		err = d.ucSession.SetValue(r.Context(), alias.SessionKey(sessionId.String()), alias.SessionValue(u.Email))
		if err != nil {
			functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignIn, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
			w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &cookie)

		csrfToken, err := GenCsrfToken(d.logger, requestId, cnst.NameHandlerSignUp, d.cfg.CsrfSecretKey, alias.SessionKey(sessionId.String()))
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
		u = sanitizer.User(u)
		uDTO := dto.NewUser(u)
		w = functions.JsonResponse(w, uDTO)
=======
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) || errors.Is(err, myerrors.BadAuthPassword) {
			w = functions.ErrorResponse(w, myerrors.BadAuthCredentialsRu, http.StatusBadRequest)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
>>>>>>> fix_csrf_test
		return
	}
	uSanitizer := sanitizer.User(u)
	uDTO := dto.NewUserData(uSanitizer)

	setCookieProps := props.GetSetCookieProps(d.ucCsrf, d.ucSession, u.Email, d.cfg.CsrfSecretKey)
	w, err = functions.SetCookie(w, r, setCookieProps)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	w = functions.JsonResponse(w, uDTO)
}

func (d *Delivery) SignOut(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	err := functions.DeleteCookiesFromDB(r, d.ucCsrf, d.ucSession)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, http.ErrNoCookie) {
			w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusForbidden)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	w, err = functions.CookieExpired(w, r)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, http.ErrNoCookie) {
			w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	w = functions.JsonResponse(w, map[string]string{"detail": "Сессия успешно завершена"})
}
