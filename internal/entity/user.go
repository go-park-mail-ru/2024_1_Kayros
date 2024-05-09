package entity

import protouser "2024_1_kayros/gen/go/user"

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

func CopyUser(u *User) *User {
	return &User{
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

func Copy(u *protouser.User) *protouser.User {
	return &protouser.User{
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

func ConvertEntityUserIntoProtoUser (u *User) *protouser.User {
	return &protouser.User{
		Id: u.Id,
		Name: u.Name,
		Phone: u.Phone,
		Email: u.Email,
		Address: u.Address,
		ImgUrl: u.ImgUrl,
		CardNumber: u.CardNumber,
		Password: u.Password,
	}
}

func ConvertProtoUserIntoEntityUser (u *protouser.User) *User {
	return &User{
		Id: u.GetId(),
		Name: u.GetName(),
		Phone: u.GetPhone(),
		Email: u.GetEmail(),
		Address: u.GetAddress(),
		ImgUrl: u.GetImgUrl(),
		CardNumber: u.GetCardNumber(),
		Password: u.GetPassword(),
	}
}

func ProtoUserIsNIL (u *protouser.User) bool {
	return u.GetId() == 0
}
