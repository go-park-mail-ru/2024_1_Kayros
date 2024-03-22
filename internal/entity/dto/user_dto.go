package dto

type User struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	ImgUrl   string `json:"img_url"`
	Password string `json:"-"`
}
