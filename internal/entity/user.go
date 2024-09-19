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
	IsVkUser   bool
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
		IsVkUser:   u.IsVkUser,
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
		IsVkUser:   u.IsVkUser,
	}
}

func ConvertEntityUserIntoProtoUser(u *User) *protouser.User {
	return &protouser.User{
		Id:         u.Id,
		Name:       u.Name,
		Phone:      u.Phone,
		Email:      u.Email,
		Address:    u.Address,
		ImgUrl:     u.ImgUrl,
		CardNumber: u.CardNumber,
		Password:   u.Password,
		IsVkUser:   u.IsVkUser,
	}
}

func ConvertProtoUserIntoEntityUser(u *protouser.User) *User {
	return &User{
		Id:         u.GetId(),
		Name:       u.GetName(),
		Phone:      u.GetPhone(),
		Email:      u.GetEmail(),
		Address:    u.GetAddress(),
		ImgUrl:     u.GetImgUrl(),
		CardNumber: u.GetCardNumber(),
		Password:   u.GetPassword(),
		IsVkUser:   u.GetIsVkUser(),
	}
}

func ProtoUserIsNIL(u *protouser.User) bool {
	return u.GetId() == 0
}
