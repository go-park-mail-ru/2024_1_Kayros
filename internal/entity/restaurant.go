package entity

type Restaurant struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImgUrl      string `json:"img_url"`
}
