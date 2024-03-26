package user

import (
	"net/http"

	"2024_1_kayros/internal/usecase/user"
)

type UserDelivery struct {
	userUsecase user.UserUsecaseInterface
}

func NewUserDelivery(uc user.UserUsecaseInterface) *UserDelivery {
	return &UserDelivery{
		userUsecase: uc,
	}
}

func (uc *UserDelivery) UserData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uc.userUsecase.GetData(w, r)
}
