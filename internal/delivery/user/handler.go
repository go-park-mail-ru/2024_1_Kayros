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

func (d *Delivery) UserAddress(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)

	address := ""
	u, err := d.ucUser.GetData(r.Context(), email)
	if err != nil && !errors.Is(err, myerrors.SqlNoRowsUserRelation) {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	unauthAddress, err := d.ucUser.GetAddressByUnauthId(r.Context(), unauthId)
	if err != nil && !errors.Is(err, myerrors.SqlNoRowsUnauthAddressRelation) {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	if unauthAddress != "" {
		address = unauthAddress
	}
	if u != nil && u.Address != "" {
		address = u.Address
	}

	address = sanitizer.Address(address)
	w = functions.JsonResponse(w, map[string]string{"address": address})
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
		w, err = functions.FlashCookie(r, w, d.ucCsrf, d.ucSession)
		if err != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
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
	data := props.GetUpdateUserDataProps(email, file, handler, u)
	uUpdated, err := d.ucUser.UpdateData(r.Context(), data)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			w, err := functions.FlashCookie(r, w, d.ucCsrf, d.ucSession)
			if err != nil {
				d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
			w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
			return
		}
		if errors.Is(err, myerrors.WrongFileExtension) {
			w = functions.ErrorResponse(w, myerrors.WrongFileExtensionRu, http.StatusBadRequest)
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

	if email != userDTO.Email {
		err = functions.DeleteCookiesFromDB(r, d.ucSession, d.ucCsrf)
		if err != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}

		setCookieProps := props.GetSetCookieProps(d.ucCsrf, d.ucSession, uUpdated.Email, d.cfg.CsrfSecretKey)
		w, err = functions.SetCookie(w, r, setCookieProps)
		if err != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		}
	}
	w = functions.JsonResponse(w, userDTO)
}

func (d *Delivery) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)

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
	isValid, err := address.Validate()
	if err != nil || !isValid {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	if email != "" {
		err = d.ucUser.UpdateAddress(r.Context(), email, address.Data)
		if err != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
				w, err = functions.FlashCookie(r, w, d.ucCsrf, d.ucSession)
				if err != nil {
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
	}
	if unauthId != "" {
		err = d.ucUser.UpdateAddressByUnauthId(r.Context(), unauthId, address.Data)
		if err != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
	}
	w = functions.JsonResponse(w, map[string]string{"detail": "Адрес успешно добавлен"})
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
			w, err = functions.FlashCookie(r, w, d.ucCsrf, d.ucSession)
			if err != nil {
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

	err = functions.DeleteCookiesFromDB(r, d.ucSession, d.ucCsrf)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	setCookieProps := props.GetSetCookieProps(d.ucCsrf, d.ucSession, email, d.cfg.CsrfSecretKey)
	w, err = functions.SetCookie(w, r, setCookieProps)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	w = functions.JsonResponse(w, map[string]string{"detail": "Пароль был успешно обновлен"})
}
