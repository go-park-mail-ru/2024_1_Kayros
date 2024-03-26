package dto

import (
	"github.com/asaskevich/govalidator"
)

type UserDTO struct {
	Id       uint64 `json:"id" valid:"-"`
	Name     string `json:"name" valid:"user_name"`
	Phone    string `json:"phone" valid:"user_phone"`
	Email    string `json:"email" valid:"user_email"`
	ImgUrl   string `json:"img_url" valid:"url"`
	Password string `json:"-" valid:"user_pwd"`
}

func (d *UserDTO) Validate() bool {
	isValid, err := govalidator.ValidateStruct(d)
	if err != nil {
		return false
	}
	return isValid
}
