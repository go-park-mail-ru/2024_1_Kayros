package entity

import userv1 "2024_1_kayros/microservices/user/proto"

type User struct {
	Id         uint64
	Name       string
	Phone      string
	Email      string
	Address    string
	ImgUrl     string
	CardNumber string
	Password   string
}

func Copy(u *userv1.User) *userv1.User {
	return &userv1.User{
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
