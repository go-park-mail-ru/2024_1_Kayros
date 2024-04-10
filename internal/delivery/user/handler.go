package user

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/delivery/auth"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/sanitizer"
	"github.com/satori/uuid"
	"go.uber.org/zap"
)

type Delivery struct {
	ucSession session.Usecase
	ucCsrf    session.Usecase
	ucUser    user.Usecase
	logger    *zap.Logger
	cfg       *config.Project
}

func NewDeliveryLayer(cfgProps *config.Project, ucSessionProps session.Usecase, ucUserProps user.Usecase, ucCsrfProps session.Usecase, loggerProps *zap.Logger) *Delivery {
	return &Delivery{
		ucUser:    ucUserProps,
		logger:    loggerProps,
		ucSession: ucSessionProps,
		ucCsrf:    ucCsrfProps,
		cfg:       cfgProps,
	}
}

func (d *Delivery) UserData(w http.ResponseWriter, r *http.Request) {
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

	u, err := d.ucUser.GetByEmail(r.Context(), email)
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUserData, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	u = sanitizer.User(u)
	uDTO := dto.NewUser(u)
	w = functions.JsonResponse(w, uDTO)

	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerUserData, cnst.DeliveryLayer)
}

func (d *Delivery) UpdateInfo(w http.ResponseWriter, r *http.Request) {
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

	// Максимальный размер фотографии 10 Mb
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		functions.LogError(d.logger, requestId, cnst.NameHandlerUploadImage, err, cnst.DeliveryLayer)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUploadImage, errors.New(myerrors.BigSizeFileError), http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BigSizeFileError, http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("img")
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			errorMsg := fmt.Sprintf("Запрос %s. Ошибка закрытия файла", requestId)
			d.logger.Error(errorMsg, zap.Error(err))
		}
	}(file)

	u := dto.GetUserFromUpdate(r)
	userUpdated, err := d.ucUser.Update(r.Context(), email, file, handler, u)
	if err != nil && strings.Contains(err.Error(), "user_email_key") {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUploadImage, errors.New(myerrors.UserAlreadyExistError), http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UserAlreadyExistError, http.StatusBadRequest)
		return
	}
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUploadImage, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	userSanitizer := sanitizer.User(userUpdated)
	userDTO := dto.NewUser(userSanitizer)

	// удалим сессии из БД
	sessionCookie, err := r.Cookie(cnst.SessionCookieName)
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

	csrfCookie, err := r.Cookie(cnst.CsrfCookieName)
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignOut, err, http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	wasDeleted, err = d.ucCsrf.DeleteKey(r.Context(), alias.SessionKey(csrfCookie.Value))
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignOut, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if !wasDeleted {
		functions.LogWarn(d.logger, requestId, cnst.NameHandlerSignOut, errors.New("Такого ключа нет в Redis"), cnst.DeliveryLayer)
	}

	sessionId := uuid.NewV4()
	expiration := time.Now().Add(14 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     cnst.SessionCookieName,
		Value:    sessionId.String(),
		Expires:  expiration,
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)

	err = d.ucSession.SetValue(r.Context(), alias.SessionKey(sessionId.String()), alias.SessionValue(userDTO.Email))
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerSignIn, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	csrfToken, err := auth.GenCsrfToken(d.logger, requestId, cnst.NameHandlerSignUp, d.cfg.CsrfSecretKey, alias.SessionKey(sessionId.String()))
	if err == nil {
		err = d.ucCsrf.SetValue(r.Context(), alias.SessionKey(csrfToken), alias.SessionValue(userDTO.Email))
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
	w = functions.JsonResponse(w, userDTO)
	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerUploadImage, cnst.DeliveryLayer)
}
