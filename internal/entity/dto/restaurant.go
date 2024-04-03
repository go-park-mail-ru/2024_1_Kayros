package dto

import (
	"2024_1_kayros/internal/entity"
	"github.com/asaskevich/govalidator"
)

// нужно будет узнать минимальную длину описания и имени
type Restaurant struct {
	Id               uint64 `json:"id" valid:"-"`
	Name             string `json:"name" valid:"-"`
	ShortDescription string `json:"short_description" valid:"-"`
	LongDescription  string `json:"long_description" valid:"-"`
	ImgUrl           string `json:"img_url" valid:"url"`
}

func (d *Restaurant) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func NewRestaurant(r *entity.Restaurant) *Restaurant {
	return &Restaurant{
		Id:               r.Id,
		Name:             r.Name,
		ShortDescription: r.ShortDescription,
		LongDescription:  r.LongDescription,
		ImgUrl:           r.ImgUrl,
	}
}

func NewRestaurantArray(restArray []*entity.Restaurant) []*Restaurant {
	restArrayDTO := make([]*Restaurant, 0, len(restArray)+1)
	for i, rest := range restArray {
		restArrayDTO[i] = NewRestaurant(rest)
	}
	return restArrayDTO
}
