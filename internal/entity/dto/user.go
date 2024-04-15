package dto

import (
	"net/http"

	"2024_1_kayros/internal/entity"
	"github.com/asaskevich/govalidator"
)

// User - DTO used for http.Request.Body JSON (update, signup, signin)
type User struct {
	Id         uint64 `json:"id" valid:"int, optional"`
	Name       string `json:"name" valid:"user_name_domain, optional"`
	Phone      string `json:"phone" valid:"user_phone_domain, optional"`
	Email      string `json:"email" valid:"user_email_domain"`
	Address    string `json:"address" valid:"user_address_domain, optional"`
	ImgUrl     string `json:"img_url" valid:"img_url_domain, optional"`
	CardNumber string `json:"card_number" valid:"user_card_number_domain, optional"`
	Password   string `json:"password" valid:"user_password_domain, optional"`
}

func (d *User) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func GetUserFromProfileForm(r *http.Request) (*entity.User, error) {
	bodyDataDTO := &User{
		Name:  r.FormValue("name"),
		Phone: r.FormValue("phone"),
		Email: r.FormValue("email"),
	}
	isValid, err := bodyDataDTO.Validate()
	if err != nil || !isValid {
		return nil, err
	}

	return &entity.User{
		Name:  bodyDataDTO.Name,
		Phone: bodyDataDTO.Phone,
		Email: bodyDataDTO.Email,
	}, nil
}

func NewUserFromSignUpForm(data *User) *entity.User {
	uDTO := &entity.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
	}
	return uDTO
}q

// UserResponse - DTO used for response to the client
type UserResponse struct {
	Id         uint64 `json:"id" valid:"int, optional"`
	Name       string `json:"name" valid:"user_name_domain"`
	Phone      string `json:"phone" valid:"user_phone_domain"`
	Email      string `json:"email" valid:"user_email_domain"`
	Address    string `json:"address" valid:"user_address_domain"`
	ImgUrl     string `json:"img_url" valid:"img_url_domain"`
	CardNumber string `json:"-" valid:"user_card_number_domain"`
	Password   string `json:"-" valid:"user_password_domain"`
}

func NewUserResponse(u *entity.User) *UserResponse {
	return &UserResponse{
		Id:         u.Id,
		Name:       u.Name,
		Phone:      u.Phone,
		Email:      u.Email,
		Address:    u.Address,
		ImgUrl:     u.ImgUrl,
		CardNumber: u.CardNumber,
		Password:   u.Password,
	}
}
