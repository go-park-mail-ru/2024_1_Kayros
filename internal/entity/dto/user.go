package dto

import (
	"2024_1_kayros/internal/entity"
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

func (d *UserDTO) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func NewUser(u *entity.User) *UserDTO {
	uDTO := &UserDTO{
		Id:       u.Id,
		Name:     u.Name,
		Phone:    u.Phone,
		Email:    u.Email,
		ImgUrl:   u.ImgUrl,
		Password: u.Password,
	}
	return uDTO
}
