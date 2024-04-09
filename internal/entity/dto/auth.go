package dto

import (
	"github.com/asaskevich/govalidator"
)

// SignIn структура данных, получаемая с формы авторизации
type SignIn struct {
	Email    string `json:"email" valid:"user_email"`
	Password string `json:"password" valid:"user_pwd"`
}

func (d *SignIn) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

// SignUp структура данных, получаемая с формы регистрации
type SignUp struct {
	Email    string `json:"email" valid:"user_email"`
	Name     string `json:"name" valid:"user_name"`
	Password string `json:"password" valid:"user_pwd"`
}

func (d *SignUp) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}
