package user

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/satori/uuid"
	"go.uber.org/zap"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/props"
	"2024_1_kayros/internal/utils/sanitizer"
<<<<<<< HEAD
)

type Delivery struct {
	ucSession       session.Usecase
	ucCsrf          session.Usecase
	ucUnauthAddress session.Usecase
	ucUser          user.Usecase
	logger          *zap.Logger
	cfg             *config.Project
=======
	"go.uber.org/zap"
)

type Delivery struct {
	ucSession session.Usecase
	ucCsrf    session.Usecase
	ucUser    user.Usecase
	cfg       *config.Project
	logger    *zap.Logger
>>>>>>> fix_csrf_test
}

func NewDeliveryLayer(cfgProps *config.Project, ucSessionProps session.Usecase, ucUserProps user.Usecase, ucCsrfProps session.Usecase, ucUnauthAddressProps session.Usecase, loggerProps *zap.Logger) *Delivery {
	return &Delivery{
<<<<<<< HEAD
		ucUser:          ucUserProps,
		logger:          loggerProps,
		ucSession:       ucSessionProps,
		ucCsrf:          ucCsrfProps,
		ucUnauthAddress: ucUnauthAddressProps,
		cfg:             cfgProps,
=======
		ucUser:    ucUserProps,
		ucSession: ucSessionProps,
		ucCsrf:    ucCsrfProps,
		cfg:       cfgProps,
		logger:    loggerProps,
>>>>>>> fix_csrf_test
	}
}

func (d *Delivery) UserData(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	u, err := d.ucUser.GetData(r.Context(), email)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if !errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
		// if error is `myerrors.RedisNoData` it's okay, because we try to delete token from database
		err = functions.DeleteCookiesFromDB(r, d.ucSession, d.ucCsrf)
		if err != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
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
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	uSanitizer := sanitizer.User(u)
	uDTO := dto.NewUserData(uSanitizer)

	w = functions.JsonResponse(w, uDTO)
}

func (d *Delivery) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	file, handler, u, err := dto.GetUpdatedUserData(r)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.BigSizeFile) {
			w = functions.ErrorResponse(w, myerrors.BigSizeFileRu, http.StatusBadRequest)
			return
		}
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		if file != nil {
			err = file.Close()
			if err != nil {
				d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			}
		}
	}(file)

<<<<<<< HEAD
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
=======
	data := props.GetUpdateUserDataProps(email, file, handler, u)
	uUpdated, err := d.ucUser.UpdateData(r.Context(), data)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.WrongFileExtension) {
			w = functions.ErrorResponse(w, myerrors.WrongFileExtensionRu, http.StatusBadRequest)
>>>>>>> fix_csrf_test
			return
		}
		if errors.Is(err, myerrors.UserAlreadyExist) {
			w = functions.ErrorResponse(w, myerrors.UserAlreadyExistRu, http.StatusBadRequest)
			return
		}
		// we don't handle `myerrors.SqlNoRowsUserRelation`, because it's internal server error
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	uSanitizer := sanitizer.User(uUpdated)
	userDTO := dto.NewUserData(uSanitizer)

	err = functions.DeleteCookiesFromDB(r, d.ucSession, d.ucCsrf)
	if err != nil {
		// we don't handle `myerrors.RedisNoData`, because it's internal server error | at first, middlewares check session_id and csrf_token
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	setCookieProps := props.GetSetCookieProps(d.ucCsrf, d.ucSession, uUpdated.Email, d.cfg.CsrfSecretKey)
	w, err = functions.SetCookie(w, r, setCookieProps)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	w = functions.JsonResponse(w, userDTO)
}

func (d *Delivery) UpdateAddress(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD
	w.Header().Set("Content-Type", "application/json")
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		err := errors.New("request_id передан не был")
		functions.LogError(d.logger, requestId, cnst.NameHandlerUpdateAddress, err, cnst.DeliveryLayer)
	} else {
		requestId = ctxRequestId.(string)
	}

	unauthToken := ""
	ctxUnauthToken := r.Context().Value(cnst.UnauthTokenCookieName)
	if ctxUnauthToken != nil {
		unauthToken = ctxUnauthToken.(string)
	}

	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	if email == "" && unauthToken == "" {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateAddress, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
=======
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
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

	var address dto.Address
	err = json.Unmarshal(body, &address)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
<<<<<<< HEAD
	if len(address.Data) < 14 || len(address.Data) > 100 {
		err = errors.New(myerrors.BadCredentialsError)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateAddress, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	if email == "" {
		err = d.ucUnauthAddress.SetValue(r.Context(), alias.SessionKey(unauthToken), alias.SessionValue(address.Data))
		if err != nil {
			functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateAddress, err, http.StatusInternalServerError, cnst.DeliveryLayer)
			w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
			return
		}
	} else {
		u, err := d.ucUser.GetByEmail(r.Context(), email)
		if err != nil || u == nil {
			err = errors.New(myerrors.InternalServerError)
			functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdateAddress, err, http.StatusInternalServerError, cnst.DeliveryLayer)
			w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
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
	}
	w = functions.JsonResponse(w, map[string]string{"address": address.Data})
	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerUpdateAddress, cnst.DeliveryLayer)
=======

	isValid, err := address.Validate()
	if address.Data == "" {
		isValid = true
	}
	if err != nil || !isValid {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	uUpdated, err := d.ucUser.UpdateAddress(r.Context(), email, address.Data)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			err = functions.DeleteCookiesFromDB(r, d.ucSession, d.ucCsrf)
			if err != nil {
				d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
			w, err = functions.CookieExpired(w, r)
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
					w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
					return
				}
			}
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	uSanitizer := sanitizer.User(uUpdated)
	uDTO := dto.NewUserData(uSanitizer)

	err = functions.DeleteCookiesFromDB(r, d.ucSession, d.ucCsrf)
	if err != nil {
		// we don't handle `myerrors.RedisNoData`, because it's internal server error | at first, middlewares check session_id and csrf_token
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	setCookieProps := props.GetSetCookieProps(d.ucCsrf, d.ucSession, email, d.cfg.CsrfSecretKey)
	w, err = functions.SetCookie(w, r, setCookieProps)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}

	w = functions.JsonResponse(w, uDTO)
>>>>>>> fix_csrf_test
}

func (d *Delivery) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
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

	var pwds dto.Passwords
	err = json.Unmarshal(body, &pwds)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	isValid, err := pwds.Validate()
	if err != nil || !isValid {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
<<<<<<< HEAD
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
=======

	data := &props.SetNewUserPasswordProps{
		Password:    pwds.Password,
		PasswordNew: pwds.PasswordNew,
	}
	err = d.ucUser.SetNewPassword(r.Context(), email, data)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.IncorrectCurrentPassword) {
			w = functions.ErrorResponse(w, myerrors.IncorrectCurrentPasswordRu, http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.NewPassword) {
			w = functions.ErrorResponse(w, myerrors.NewPasswordRu, http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			err = functions.DeleteCookiesFromDB(r, d.ucSession, d.ucCsrf)
			if err != nil {
				d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				if errors.Is(err, myerrors.RedisNoData) {
					w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
					return
				}
			}
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
>>>>>>> fix_csrf_test
		return
	}

	w = functions.JsonResponse(w, map[string]string{"detail": "Пароль был успешно обновлен"})
<<<<<<< HEAD
	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerUpdatePassword, cnst.DeliveryLayer)
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
=======
>>>>>>> fix_csrf_test
}
