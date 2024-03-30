package dto

// нужно будет узнать минимальную длину описания и имени

type RestaurantDTO struct {
	Id               uint64 `json:"id" valid:"-"`
	Name             string `json:"name" valid:"-"`
	ShortDescription string `json:"short_description,omitempty" valid:"-"`
	LongDescription  string `json:"long_description,omitempty" valid:"-"`
	ImgUrl           string `json:"img_url" valid:"url"`
}
