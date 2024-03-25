package dto

type RestaurantDTO struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImgUrl      string `json:"img_url"`
}
