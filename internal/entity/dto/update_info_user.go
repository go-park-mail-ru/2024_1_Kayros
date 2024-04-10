package dto

import (
	"net/http"

	"2024_1_kayros/internal/entity"
	"github.com/asaskevich/govalidator"
)

type UserUpdate struct {
	Name  string `json:"name" valid:"user_name"`
	Phone string `json:"phone" valid:"user_phone"`
	Email string `json:"email" valid:"user_email"`
}

func (d *UserUpdate) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func GetUserFromUpdate(r *http.Request) *entity.User {
	bodyDataDTO := &UserUpdate{
		Name:  r.FormValue("name"),
		Phone: r.FormValue("phone"),
		Email: r.FormValue("email"),
	}
	isValid, err := bodyDataDTO.Validate()
	if err != nil || !isValid {
		return nil
	}

	return &entity.User{
		Name:  bodyDataDTO.Name,
		Phone: bodyDataDTO.Phone,
		Email: bodyDataDTO.Email,
	}
}
