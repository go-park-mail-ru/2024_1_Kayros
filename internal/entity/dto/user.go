package dto

import (
	"mime/multipart"
	"net/http"

	"2024_1_kayros/internal/entity"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/myerrors"
	"github.com/asaskevich/govalidator"
)

// UserUpdate - DTO used for unmarshalling http.Request.Body in format JSON (need for updating user data from profile)
type UserUpdate struct {
	Name  string `json:"name" valid:"user_name_domain"`
	Phone string `json:"phone" valid:"user_phone_domain, optional"`
	Email string `json:"email" valid:"user_email_domain"`
}

func (d *UserUpdate) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func GetUpdatedUserData(r *http.Request) (multipart.File, *multipart.FileHeader, *entity.User, error) {
	bodyDataDTO := &UserUpdate{
		Name:  r.FormValue("name"),
		Phone: r.FormValue("phone"),
		Email: r.FormValue("email"),
	}
	isValid, err := bodyDataDTO.Validate()
	if err != nil || !isValid {
		return nil, nil, nil, err
	}
	u := &entity.User{
		Name:  bodyDataDTO.Name,
		Phone: bodyDataDTO.Phone,
		Email: bodyDataDTO.Email,
	}

	file, handler, err := r.FormFile("img")
	if err != nil {
		return nil, nil, u, err
	}
	if handler.Size > cnst.UploadedFileMaxSize {
		return file, handler, u, myerrors.BigSizeFile
	}

	return file, handler, u, nil
}

type UserSignUp struct {
	Name     string `json:"name" valid:"user_name_domain"`
	Email    string `json:"email" valid:"user_email_domain"`
	Password string `json:"password" valid:"user_password_domain"`
}

func NewUserFromSignUpForm(data *UserSignUp) *entity.User {
	uDTO := &entity.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
	}
	return uDTO
}

// UserGet - DTO used for handler 'UserData' method GET
type UserGet struct {
	Id      uint64 `json:"id" valid:"int, optional"`
	Name    string `json:"name" valid:"user_name_domain"`
	Phone   string `json:"phone" valid:"user_phone_domain"`
	Email   string `json:"email" valid:"user_email_domain"`
	Address string `json:"address" valid:"user_address_domain"`
	ImgUrl  string `json:"img_url" valid:"img_url_domain"`
}

// NewUserData - function used to form response for receiving detailed information about user
func NewUserData(u *entity.User) *UserGet {
	return &UserGet{
		Id:      u.Id,
		Name:    u.Name,
		Phone:   u.Phone,
		Email:   u.Email,
		Address: u.Address,
		ImgUrl:  u.ImgUrl,
	}
}

// Address - DTO used for handler 'UpdateAddress' method PUT
type Address struct {
	Data string `json:"user_address_domain"`
}

func (d *Address) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

type Passwords struct {
	Password    string `json:"password" valid:"user_password_domain"`
	PasswordNew string `json:"new_password" valid:"user_password_domain"`
}

func (d *Passwords) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}
