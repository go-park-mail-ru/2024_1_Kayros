package user

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/props"
	"2024_1_kayros/internal/utils/sanitizer"

	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type Delivery struct {
	ucSession session.Usecase
	ucUser    user.Usecase
	cfg       *config.ProjectConfiguration
	logger    *zap.Logger
	metrics   *metrics.Metrics
}

func NewDeliveryLayer(ucSessionProps session.Usecase, ucUserProps user.Usecase, loggerProps *zap.Logger, metrics *metrics.Metrics) *Delivery {
	return &Delivery{
		ucUser:    ucUserProps,
		ucSession: ucSessionProps,
		cfg:       &config.Config,
		logger:    loggerProps,
		metrics:   metrics,
	}
}

func (d *Delivery) UserAddress(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)

	isUserAddress := r.URL.Query().Get("user_address")
	isUserAddress = strings.TrimSpace(isUserAddress)

	if isUserAddress == "true" && email == "" {
		d.logger.Error("unauthorized user can't get authorized user's address", zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadRequestGetAddress, http.StatusBadRequest)
		return
	}

	address, err := d.ucUser.UserAddress(r.Context(), email, unauthId, isUserAddress)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	address = sanitizer.Address(address)

	functions.JsonResponse(w, &dto.Address{Data: address})
}

func (d *Delivery) UpdateUnauthAddress(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)
	if email == "" {
		functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}
	if unauthId == "" {
		functions.ErrorResponse(w, myerrors.BadRequestUpdateUnauthAddress, http.StatusBadRequest)
		return
	}

	u, err := d.ucUser.GetData(r.Context(), email)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	if u != nil {
		err = d.ucUser.UpdateUnauthAddress(r.Context(), u.Address, unauthId)
		if err != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
		functions.JsonResponse(w, &dto.ResponseDetail{Detail: "Адрес успешно выбран"})
	}
	functions.ErrorResponse(w, errors.New("Адрес не был выбран"), http.StatusInternalServerError)
}

func (d *Delivery) UserData(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	u, err := d.ucUser.GetData(r.Context(), email)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if !errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
		w, err = functions.FlashCookie(r, w, d.ucSession, &d.cfg.Redis, d.metrics)
		if err != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
		functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}
	uDTO := dto.NewUserData(sanitizer.User(u))
	functions.JsonResponse(w, uDTO)
}

func (d *Delivery) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	file, handler, u, err := dto.GetUpdatedUserData(r)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.BigSizeFile) {
			functions.ErrorResponse(w, myerrors.BigSizeFileRu, http.StatusBadRequest)
			return
		}
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
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
	uUpdated, err := d.ucUser.UpdateData(r.Context(), props.GetUpdateUserDataProps(email, file, handler, u))
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			w, err := functions.FlashCookie(r, w, d.ucSession, &d.cfg.Redis, d.metrics)
			if err != nil {
				d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
			functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
			return
		}
		if errors.Is(err, myerrors.WrongFileExtension) {
			functions.ErrorResponse(w, myerrors.WrongFileExtensionRu, http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.UserAlreadyExist) {
			functions.ErrorResponse(w, myerrors.UserAlreadyExistRu, http.StatusBadRequest)
			return
		}
		// we don't handle `myerrors.SqlNoRowsUserRelation`, because it's internal server error
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	userDTO := dto.NewUserData(sanitizer.User(uUpdated))

	if email != userDTO.Email {
		err = functions.DeleteCookiesFromDB(r, d.ucSession, &d.cfg.Redis, d.metrics)
		if err != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}

		w, err = functions.SetCookie(w, r, d.ucSession, uUpdated.Email, d.cfg)
		if err != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		}
	}
	functions.JsonResponse(w, userDTO)
}

func (d *Delivery) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)

	isUserAddress := r.URL.Query().Get("user_address")
	isUserAddress = strings.TrimSpace(isUserAddress)

	if isUserAddress == "true" && email == "" {
		d.logger.Error("unauthorized user can't update authorized user's email", zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadRequestUpdateAddress, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	var address dto.Address
	err = easyjson.Unmarshal(body, &address)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	isValid, err := address.Validate()
	if err != nil || !isValid {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	err = d.ucUser.UpdateAddress(r.Context(), email, unauthId, address.Data, isUserAddress)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			w, err = functions.FlashCookie(r, w, d.ucSession, &d.cfg.Redis, d.metrics)
			if err != nil {
				d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
			functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
			return
		}
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	functions.JsonResponse(w, &dto.ResponseDetail{Detail: "Адрес успешно добавлен"})
}

func (d *Delivery) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	var pwds dto.Passwords
	err = easyjson.Unmarshal(body, &pwds)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	isValid, err := pwds.Validate()
	if err != nil || !isValid {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	err = d.ucUser.SetNewPassword(r.Context(), email, pwds.Password, pwds.PasswordNew)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.IncorrectCurrentPassword) {
			functions.ErrorResponse(w, myerrors.IncorrectCurrentPasswordRu, http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.NewPassword) {
			functions.ErrorResponse(w, myerrors.NewPasswordRu, http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			w, err = functions.FlashCookie(r, w, d.ucSession, &d.cfg.Redis, d.metrics)
			if err != nil {
				d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
			functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
			return
		}
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	err = functions.DeleteCookiesFromDB(r, d.ucSession, &d.cfg.Redis, d.metrics)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	w, err = functions.SetCookie(w, r, d.ucSession, email, d.cfg)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	functions.JsonResponse(w, &dto.ResponseDetail{Detail: "Пароль был успешно обновлен"})
}
