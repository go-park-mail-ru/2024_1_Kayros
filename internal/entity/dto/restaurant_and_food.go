package dto

import (
	"fmt"

	"github.com/asaskevich/govalidator"

	"2024_1_kayros/internal/entity"
)

type RestaurantAndFood struct {
	Id               uint64         `json:"id" valid:"-"`
	Name             string         `json:"name" valid:"-"`
	ShortDescription string         `json:"short_description" valid:"-"`
	LongDescription  string         `json:"long_description" valid:"-"`
	ImgUrl           string         `json:"img_url" valid:"url"`
	Food             []*entity.Food `json:"food" valid:"-"`
}

func (d *RestaurantAndFood) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func NewRestaurantAndFood(r *entity.Restaurant, foodArray []*entity.Food) *RestaurantAndFood {
	fmt.Println("dto")
	return &RestaurantAndFood{
		Id:               r.Id,
		Name:             r.Name,
		ShortDescription: r.ShortDescription,
		LongDescription:  r.LongDescription,
		ImgUrl:           r.ImgUrl,
		Food:             foodArray,
	}
}
