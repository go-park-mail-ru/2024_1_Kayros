package user

import (
	"log"
	"mime/multipart"
	"net/http"

	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/user"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
)

type Delivery struct {
	ucUser user.Usecase
}

func NewDeliveryLayer(ucUserProps user.Usecase) *Delivery {
	return &Delivery{
		ucUser: ucUserProps,
	}
}

func (d *Delivery) UserData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := r.Context().Value("email")
	if email == nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	u, err := d.ucUser.GetByEmail(r.Context(), email.(string))
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	uDTO := dto.NewUser(u)
	w = functions.JsonResponse(w, uDTO)
}

func (d *Delivery) UpdateImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := r.Context().Value("email")
	if email == nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	// Максимальный размер фотографии 10 Mb
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.BigSizeFileError, http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("img")
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(file)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}

	err = d.ucUser.UploadImageByEmail(r.Context(), file, handler, email.(string))
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	u, err := d.ucUser.GetByEmail(r.Context(), email.(string))
	userDTO := dto.NewUser(u)
	w = functions.JsonResponse(w, userDTO)
}
