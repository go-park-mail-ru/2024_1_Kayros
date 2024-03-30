package functions

import (
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/entity/dto"
)

func ConvIntoUserDTO(u *entity.User) *dto.UserDTO {
	uDTO := &dto.UserDTO{
		Id:       u.Id,
		Name:     u.Name,
		Phone:    u.Phone,
		Email:    u.Email,
		ImgUrl:   u.ImgUrl,
		Password: u.Password,
	}
	return uDTO
}
