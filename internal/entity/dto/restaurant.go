package dto

import (
	"2024_1_kayros/internal/entity"
	"github.com/asaskevich/govalidator"
)

// нужно будет узнать минимальную длину описания и имени

type RestaurantDTO struct {
	Id               uint64 `json:"id" valid:"-"`
	Name             string `json:"name" valid:"-"`
	ShortDescription string `json:"short_description" valid:"-"`
	LongDescription  string `json:"long_description" valid:"-"`
	ImgUrl           string `json:"img_url" valid:"url"`
}

func (d *RestaurantDTO) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func NewRestaurant(r *entity.Restaurant) *RestaurantDTO {
	return &RestaurantDTO{
		Id:               r.Id,
		Name:             r.Name,
		ShortDescription: r.ShortDescription,
		LongDescription:  r.LongDescription,
		ImgUrl:           r.ImgUrl,
	}
}