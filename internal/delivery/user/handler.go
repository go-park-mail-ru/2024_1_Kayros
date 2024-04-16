package user

import (
	"encoding/json"
	"errors"
	"io"
	"log"
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
	"github.com/asaskevich/govalidator"
	"github.com/satori/uuid"
)

type Delivery struct {
	ucSession session.Usecase
	ucCsrf    session.Usecase
	ucUser    user.Usecase
	cfg       *config.Project
}

func NewDeliveryLayer(cfgProps *config.Project, ucSessionProps session.Usecase, ucUserProps user.Usecase, ucCsrfProps session.Usecase) *Delivery {
	return &Delivery{
		ucUser:    ucUserProps,
		ucSession: ucSessionProps,
		ucCsrf:    ucCsrfProps,
		cfg:       cfgProps,
	}
}

func (d *Delivery) UserData(w http.ResponseWriter, r *http.Request) {
	logger, err := functions.GetCtxLogger(r)
	if err != nil {
		return
	}
	requestId, err := functions.GetCtxRequestId(r)
	if err != nil {
		return
	}
	email, err := functions.GetCtxEmail(r)
	if err != nil {
		return
	}

	u, err := d.ucUser.GetUserData(r.Context(), email, requestId, logger)
	if err != nil {
		// нужно обработать кастомные ошибки БД и тут их обрабатывать
		return
	}

	uSanitizer := sanitizer.User(u)
	uDTO := dto.NewUserData(uSanitizer)
	w = functions.JsonResponse(w, uDTO)
}

func (d *Delivery) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	logger, err := functions.GetCtxLogger(r)
	if err != nil {
		return
	}
	requestId, err := functions.GetCtxRequestId(r)
	if err != nil {
		return
	}
	email, err := functions.GetCtxEmail(r)
	if err != nil {
		return
	}

	file, handler, u, err := dto.GetUpdatedUserData(r)
	defer func(file multipart.File) {
		if file != nil {
			err = file.Close()
			if err != nil {
				log.Println("Error of closing file")
			}
		}
	}(file)
	if err != nil {
		if strings.Contains(err.Error(), "email") {
			err = errors.New("Некорректный пароль")
			return
		}
		if strings.Contains(err.Error(), "phone") {
			err = errors.New("Некорректный номер телефона")
			return
		}
		if strings.Contains(err.Error(), "name") {
			err = errors.New("Некорректное имя")
			return
		}
		err = errors.New("Некорректные данные")
		return
	}

	uUpdated, err := d.ucUser.UpdateUserData(r.Context(), email, file, handler, u, requestId, logger)
	if err != nil {
		return
	}
	uSanitizer := sanitizer.User(uUpdated)
	userDTO := dto.NewUserData(uSanitizer)

	err = functions.DeleteCookies(r, d.ucSession, d.ucCsrf)
	if err != nil {
		return
	}
	functions.SetCookie(w, r, d.ucCsrf, d.ucSession, email, d.cfg.CsrfSecretKey)
	w = functions.JsonResponse(w, userDTO)
}

func (d *Delivery) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	logger, err := functions.GetCtxLogger(r)
	if err != nil {
		return
	}
	requestId, err := functions.GetCtxRequestId(r)
	if err != nil {
		return
	}
	email, err := functions.GetCtxEmail(r)
	if err != nil {
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}

	var address addressData
	err = json.Unmarshal(body, &address)
	if err != nil {
		return
	}
	isValid, err := address.Validate()
	if err != nil {
		return
	}
	if !isValid {
		return
	}
	uUpdated, err := d.ucUser.UpdateUserAddress(r.Context(), email, address.Data, requestId, logger)
	if err != nil {
		return
	}

	uSanitizer := sanitizer.User(uUpdated)
	uDTO := dto.NewUserData(uSanitizer)
	err = functions.DeleteCookies(r, d.ucSession, d.ucCsrf)
	if err != nil {
		return
	}
	functions.SetCookie(w, r, d.ucCsrf, d.ucSession, email, d.cfg.CsrfSecretKey)

	w = functions.JsonResponse(w, uDTO)
}

func (d *Delivery) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	logger, err := functions.GetCtxLogger(r)
	if err != nil {
		return
	}
	requestId, err := functions.GetCtxRequestId(r)
	if err != nil {
		return
	}
	email, err := functions.GetCtxEmail(r)
	if err != nil {
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}

	var password passwordData
	err = json.Unmarshal(body, &password)
	if err != nil {
		return
	}

	isValid, err := password.Validate()
	if err != nil || !isValid {
		err = errors.New(myerrors.BadCredentialsError)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdatePassword, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}

	isEqual, err := d.ucUser.CheckPassword(r.Context(), email, password.Data)
	if err != nil {
		err = errors.New(myerrors.InternalServerError)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdatePassword, err, http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if isEqual {
		err = errors.New(myerrors.EqualPasswordsError)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdatePassword, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.EqualPasswordsError, http.StatusBadRequest)
		return
	}
	_, err = d.ucUser.SetNewPassword(r.Context(), email, password.Data)
	if err != nil {
		err = errors.New(myerrors.BadCredentialsError)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdatePassword, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	//

	csrfCookie, err := r.Cookie(cnst.CsrfCookieName)
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, err, http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	wasDeleted, err := d.ucCsrf.DeleteKey(r.Context(), alias.SessionKey(csrfCookie.Value))
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

	err = d.ucSession.SetValue(r.Context(), alias.SessionKey(sessionId.String()), alias.SessionValue(uUpdated.Email))
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	csrfToken, err := auth.GenCsrfToken(d.logger, requestId, cnst.NameHandlerUpdateUser, d.cfg.CsrfSecretKey, alias.SessionKey(sessionId.String()))
	if err == nil {
		err = d.ucCsrf.SetValue(r.Context(), alias.SessionKey(csrfToken), alias.SessionValue(uUpdated.Email))
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

	w = functions.JsonResponse(w, map[string]string{"detail": "Пароль был успешно обновлен"})
	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, cnst.DeliveryLayer)
}

type addressData struct {
	Data string `json:"user_address_domain"`
}

func (d *addressData) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

type passwordData struct {
	Data string `json:"password" valid:"user_pwd"`
}

func (d *passwordData) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}
