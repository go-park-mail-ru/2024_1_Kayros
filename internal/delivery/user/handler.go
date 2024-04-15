package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/satori/uuid"
	"go.uber.org/zap"

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
		functions.LogError(d.logger, requestId, cnst.NameHandlerUpdateUser, err, cnst.DeliveryLayer)
	} else {
		requestId = ctxRequestId.(string)
	}

	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	if email == "" {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	// Максимальный размер фотографии 10 Mb
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		functions.LogError(d.logger, requestId, cnst.NameHandlerUpdateUser, err, cnst.DeliveryLayer)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, errors.New(myerrors.BigSizeFileError), http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BigSizeFileError, http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("img")
	defer func(file multipart.File) {
		if file != nil {
			err := file.Close()
			if err != nil {
				errorMsg := fmt.Sprintf("Запрос %s. Ошибка закрытия файла", requestId)
				d.logger.Error(errorMsg, zap.Error(err))
			}
		}
	}(file)

	u := dto.GetUserFromUpdate(r)
	if u == nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, errors.New(myerrors.BadCredentialsError), http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	userUpdated, err := d.ucUser.Update(r.Context(), email, file, handler, u)
	if err != nil {
		if strings.Contains(err.Error(), "user_email_key") {
			functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, errors.New(myerrors.UserAlreadyExistError), http.StatusBadRequest, cnst.DeliveryLayer)
			w = functions.ErrorResponse(w, myerrors.UserAlreadyExistError, http.StatusBadRequest)
			return
		}
		if strings.Contains(err.Error(), "не может быть") {
			functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, err, http.StatusBadRequest, cnst.DeliveryLayer)
			w = functions.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	userSanitizer := sanitizer.User(userUpdated)
	userDTO := dto.NewUser(userSanitizer)

	// удалим сессии из БД
	sessionCookie, err := r.Cookie(cnst.SessionCookieName)
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, err, http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	wasDeleted, err := d.ucSession.DeleteKey(r.Context(), alias.SessionKey(sessionCookie.Value))
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if !wasDeleted {
		functions.LogWarn(d.logger, requestId, cnst.NameHandlerUpdateUser, errors.New("Такого ключа нет в Redis"), cnst.DeliveryLayer)
	}

	csrfCookie, err := r.Cookie(cnst.CsrfCookieName)
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, err, http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	wasDeleted, err = d.ucCsrf.DeleteKey(r.Context(), alias.SessionKey(csrfCookie.Value))
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if !wasDeleted {
		functions.LogWarn(d.logger, requestId, cnst.NameHandlerUpdateUser, errors.New("Такого ключа нет в Redis"), cnst.DeliveryLayer)
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
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	csrfToken, err := auth.GenCsrfToken(d.logger, requestId, cnst.NameHandlerUpdateUser, d.cfg.CsrfSecretKey, alias.SessionKey(sessionId.String()))
	if err == nil {
		err = d.ucCsrf.SetValue(r.Context(), alias.SessionKey(csrfToken), alias.SessionValue(userDTO.Email))
		if err != nil {
			functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
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
	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, cnst.DeliveryLayer)
}

func (d *Delivery) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		err := errors.New("request_id передан не был")
		functions.LogError(d.logger, requestId, cnst.NameHandlerUpdateAddress, err, cnst.DeliveryLayer)
	} else {
		requestId = ctxRequestId.(string)
	}

	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	if email == "" {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateAddress, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateAddress, err, http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	var address addressData
	err = json.Unmarshal(body, &address)
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateAddress, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	u, err := d.ucUser.GetByEmail(r.Context(), email)
	if err != nil || u == nil {
		err = errors.New(myerrors.InternalServerError)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateAddress, err, http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if len(address.Data) < 14 || len(address.Data) > 100 {
		err = errors.New(myerrors.BadCredentialsError)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateAddress, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	u.Address = address.Data
	uUpdated, err := d.ucUser.Update(r.Context(), email, nil, nil, u)
	if err != nil || uUpdated == nil {
		err = errors.New(myerrors.InternalServerError)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateAddress, err, http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	u = sanitizer.User(uUpdated)
	w = functions.JsonResponse(w, u)
	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, cnst.DeliveryLayer)
}

func (d *Delivery) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		err := errors.New("request_id передан не был")
		functions.LogError(d.logger, requestId, cnst.NameHandlerUpdateAddress, err, cnst.DeliveryLayer)
	} else {
		requestId = ctxRequestId.(string)
	}

	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	if email == "" {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateAddress, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateAddress, err, http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	var password passwordData
	err = json.Unmarshal(body, &password)
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateAddress, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}

	isValid, err := password.Validate()
	if err != nil || !isValid {
		err = errors.New(myerrors.BadCredentialsError)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdatePassword, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	// сравниваем старый пароль с тем, что в базе
	isEqual, err := d.ucUser.CheckPassword(r.Context(), email, password.Password)
	if err != nil {
		err = errors.New(myerrors.InternalServerError)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdatePassword, err, http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	// если они не совпадают
	if !isEqual {
		err = errors.New(myerrors.WrongPasswordError)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdatePassword, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.WrongPasswordError, http.StatusBadRequest)
		return
	}
	// они совпадают, значит мы можем поменять пароль пользователю
	// проверяем, что старый и новый пароль должны быть разными
	if password.Password == password.NewPassword {
		err = errors.New(myerrors.EqualPasswordsError)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdatePassword, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.EqualPasswordsError, http.StatusBadRequest)
		return
	}
	_, err = d.ucUser.SetNewPassword(r.Context(), email, password.NewPassword)
	if err != nil {
		err = errors.New(myerrors.BadCredentialsError)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdatePassword, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	w = functions.JsonResponse(w, map[string]string{"detail": "Пароль был успешно обновлен"})
	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, cnst.DeliveryLayer)
}

type addressData struct {
	Data string `json:"address"`
}

type passwordData struct {
	Password    string `json:"password" valid:"user_pwd"`
	NewPassword string `json:"new_password" valid:"user_pwd"`
}

func (d *passwordData) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}
