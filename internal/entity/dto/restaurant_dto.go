package dto

// нужно будет узнать минимальную длину описания и имени
type RestaurantDTO struct {
	Id          uint64 `json:"id" valid:"-"`
	Name        string `json:"name" valid:"-"`
	Description string `json:"description" valid:"-"`
	ImgUrl      string `json:"img_url" valid:"url"`
}
