package dto

import "github.com/asaskevich/govalidator"

// SignInDTO структура данных, получаемая с формы авторизации
type SignInDTO struct {
	Email    string `json:"email" valid:"user_email"`
	Password string `json:"-" valid:"user_pwd"`
}

func (d *SignInDTO) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

// SignUpDTO структура данных, получаемая с формы регистрации
type SignUpDTO struct {
	Email    string `json:"email" valid:"user_email"`
	Name     string `json:"name" valid:"user_name"`
	Password string `json:"-" valid:"user_pwd"`
}

func (d *SignUpDTO) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}
