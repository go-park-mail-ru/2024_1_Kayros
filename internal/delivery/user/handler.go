package user

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/user"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"go.uber.org/zap"
)

type Delivery struct {
	ucUser user.Usecase
	logger *zap.Logger
}

func NewDeliveryLayer(ucUserProps user.Usecase, loggerProps *zap.Logger) *Delivery {
	return &Delivery{
		ucUser: ucUserProps,
		logger: loggerProps,
	}
}

func (d *Delivery) UserData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := r.Context().Value("request_id").(string)
	email := r.Context().Value("email")
	if email == nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUserData, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	u, err := d.ucUser.GetByEmail(r.Context(), email.(string))
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUserData, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	uDTO := dto.NewUser(u)
	w = functions.JsonResponse(w, uDTO)

	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerUserData, cnst.DeliveryLayer)
}

func (d *Delivery) UploadImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := r.Context().Value("request_id").(string)
	email := r.Context().Value("email")
	if email == nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUploadImage, errors.New(myerrors.UnauthorizedError), http.StatusUnauthorized, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	// Максимальный размер фотографии 10 Mb
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
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
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUploadImage, errors.New(myerrors.BadCredentialsError), http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}

	err = d.ucUser.UploadImageByEmail(r.Context(), file, handler, email.(string))
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUploadImage, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	u, err := d.ucUser.GetByEmail(r.Context(), email.(string))
	if err != nil {
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUploadImage, errors.New(myerrors.InternalServerError), http.StatusInternalServerError, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	userDTO := dto.NewUser(u)
	w = functions.JsonResponse(w, userDTO)

	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerUploadImage, cnst.DeliveryLayer)
}
