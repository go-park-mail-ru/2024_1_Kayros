package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/auth"
	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/sanitizer"
	"go.uber.org/zap"
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
	requestId, err := functions.GetCtxRequestId(r)
	if err != nil {
		d.logger.Error(err.Error())
	}
	email, err := functions.GetCtxEmail(r)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
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

	uAuth, err := d.ucAuth.SignUpUser(r.Context(), email, u)
	if err != nil {
		if errors.Is(err, myerrors.UserAlreadyExist) {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.UserAlreadyExistRu, http.StatusBadRequest)
			return
		}
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	u = sanitizer.User(uAuth)
	uDTO := dto.NewUserData(u)

	w, err = functions.SetCookie(w, r, d.ucCsrf, d.ucSession, email, d.cfg.CsrfSecretKey)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	w = functions.JsonResponse(w, uDTO)
}

func (d *Delivery) SignIn(w http.ResponseWriter, r *http.Request) {
	requestId, err := functions.GetCtxRequestId(r)
	if err != nil {
		d.logger.Error(err.Error())
	}
	email, err := functions.GetCtxEmail(r)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
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
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) || errors.Is(err, myerrors.BadAuthCredentials) {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.BadAuthCredentialsRu, http.StatusBadRequest)
			return
		}
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	uSanitizer := sanitizer.User(u)
	uDTO := dto.NewUserData(uSanitizer)

	w, err = functions.SetCookie(w, r, d.ucCsrf, d.ucSession, email, d.cfg.CsrfSecretKey)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	w = functions.JsonResponse(w, uDTO)
	return
}

func (d *Delivery) SignOut(w http.ResponseWriter, r *http.Request) {
	requestId, err := functions.GetCtxRequestId(r)
	if err != nil {
		d.logger.Error(err.Error())
	}
	_, err = functions.GetCtxEmail(r)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	sessionCookie, errSession := r.Cookie(cnst.SessionCookieName)
	csrfCookie, errCsrf := r.Cookie(cnst.CsrfCookieName)

	err = functions.DeleteCookiesFromDB(r, d.ucCsrf, d.ucSession)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) && errSession != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
			return
		}
		if errors.Is(err, http.ErrNoCookie) && errCsrf != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusForbidden)
			return
		}
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}

	sessionCookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, sessionCookie)

	csrfCookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, csrfCookie)

	w = functions.JsonResponse(w, map[string]string{"detail": "Сессия успешно завершена"})
}

func GenCsrfToken(secretKey string, sessionId alias.SessionKey) (string, error) {
	// Создание csrf_token
	csrfToken, err := functions.HashCsrf(secretKey, string(sessionId))
	if err != nil {
		return "", err
	}
	return csrfToken, nil
}
