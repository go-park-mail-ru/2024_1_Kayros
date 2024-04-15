package entity

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

func Copy(u *User) *User {
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
