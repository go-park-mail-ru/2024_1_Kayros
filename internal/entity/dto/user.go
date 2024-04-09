package dto

import (
	"2024_1_kayros/internal/entity"
	"github.com/asaskevich/govalidator"
)

type User struct {
	Id         uint64 `json:"id" valid:"-"`
	Name       string `json:"name" valid:"user_name"`
	Phone      string `json:"phone" valid:"user_phone"`
	Email      string `json:"email" valid:"user_email"`
	Address    string `json:"address" valid:"-"` // нужно добавить валидацию для адреса
	ImgUrl     string `json:"img_url" valid:"url"`
	CardNumber string `json:"-" valid:"user_card_number, omitempty"`
	Password   string `json:"-" valid:"user_pwd, omitempty"`
}

func (d *User) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func NewUser(u *entity.User) *User {
	uDTO := &User{
		Id:       u.Id,
		Name:     u.Name,
		Phone:    u.Phone,
		Email:    u.Email,
		ImgUrl:   u.ImgUrl,
		Password: u.Password,
	}
	return uDTO
}

func NewUserFromSignUp(data *SignUp) *entity.User {
	uDTO := &entity.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
	}
	return uDTO
}
