package user

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/props"
	"2024_1_kayros/internal/utils/sanitizer"
	"go.uber.org/zap"
)

type Delivery struct {
	ucSession session.Usecase
	ucCsrf    session.Usecase
	ucUser    user.Usecase
	cfg       *config.Project
	logger    *zap.Logger
}

func NewDeliveryLayer(cfgProps *config.Project, ucSessionProps session.Usecase, ucUserProps user.Usecase, ucCsrfProps session.Usecase, loggerProps *zap.Logger) *Delivery {
	return &Delivery{
		ucUser:    ucUserProps,
		ucSession: ucSessionProps,
		ucCsrf:    ucCsrfProps,
		cfg:       cfgProps,
		logger:    loggerProps,
	}
}

func (d *Delivery) UserData(w http.ResponseWriter, r *http.Request) {
	requestId, err := functions.GetCtxRequestId(r)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	email, err := functions.GetCtxEmail(r)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	u, err := d.ucUser.GetData(r.Context(), email)
	if err != nil {
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			err = functions.DeleteCookiesFromDB(r, d.ucSession, d.ucCsrf)
			if err != nil {
				if errors.Is(err, myerrors.RedisNoData) {
					d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
					w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
					return
				}
				d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
			w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	uSanitizer := sanitizer.User(u)
	uDTO := dto.NewUserData(uSanitizer)
	w = functions.JsonResponse(w, uDTO)
}

func (d *Delivery) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	requestId, err := functions.GetCtxRequestId(r)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		return
	}
	email, err := functions.GetCtxEmail(r)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	// defer is declared immediately after the method because if parsing is successful but the file is large,
	// we return the file
	file, handler, u, err := dto.GetUpdatedUserData(r)
	defer func(file multipart.File) {
		if file != nil {
			err = file.Close()
			if err != nil {
				d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			}
		}
	}(file)
	if err != nil {
		if errors.Is(err, myerrors.BigSizeFile) {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.BigSizeFileRu, http.StatusBadRequest)
			return
		}
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	data := &props.UpdateUserDataProps{
		Email:           email,
		File:            file,
		Handler:         handler,
		UserPropsUpdate: u,
	}
	uUpdated, err := d.ucUser.UpdateData(r.Context(), data)
	if err != nil {
		if errors.Is(err, myerrors.WrongFileExtension) {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.WrongFileExtensionRu, http.StatusBadRequest)
			return
		}
		// we don't handle `myerrors.SqlNoRowsUserRelation`, because it's internal server error
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
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
	w, err = functions.SetCookie(w, r, d.ucCsrf, d.ucSession, email, d.cfg.CsrfSecretKey)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	w = functions.JsonResponse(w, userDTO)
}

func (d *Delivery) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	requestId, err := functions.GetCtxRequestId(r)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		return
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
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	var address dto.Address
	err = json.Unmarshal(body, &address)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	isValid, err := address.Validate()
	if err != nil || !isValid {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	uUpdated, err := d.ucUser.UpdateAddress(r.Context(), email, address.Data)
	if err != nil {
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
			return
		}
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
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
	w, err = functions.SetCookie(w, r, d.ucCsrf, d.ucSession, email, d.cfg.CsrfSecretKey)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}

	w = functions.JsonResponse(w, uDTO)
}

func (d *Delivery) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	requestId, err := functions.GetCtxRequestId(r)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		return
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
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
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
	data := &props.SetNewUserPasswordProps{
		Password:    pwds.Password,
		PasswordNew: pwds.PasswordNew,
	}
	err = d.ucUser.SetNewPassword(r.Context(), email, data)
	if err != nil {
		if errors.Is(err, myerrors.Password) {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.PasswordRu, http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.NewPasswordRu) {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.NewPasswordRu, http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			err = functions.DeleteCookiesFromDB(r, d.ucSession, d.ucCsrf)
			if err != nil {
				if errors.Is(err, myerrors.RedisNoData) {
					d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
					w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
					return
				}
				d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
		}
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	w = functions.JsonResponse(w, map[string]string{"detail": "Пароль был успешно обновлен"})
}
