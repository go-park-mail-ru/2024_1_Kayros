package user

import (
	"net/http"

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

	uDTO := functions.ConvIntoUserDTO(u)

	w = functions.JsonResponse(w, uDTO)
}
