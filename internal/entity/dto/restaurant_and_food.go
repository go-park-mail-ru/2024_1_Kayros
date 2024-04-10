package dto

import (
	"github.com/asaskevich/govalidator"

	"2024_1_kayros/internal/entity"
)

type Category struct {
	Id   uint64  `json:"id" valid:"-"`
	Name string  `json:"name" valid:"-"`
	Food []*Food `json:"food" valid:"-"`
}

type RestaurantAndFood struct {
	Id               uint64      `json:"id" valid:"-"`
	Name             string      `json:"name" valid:"-"`
	ShortDescription string      `json:"short_description" valid:"-"`
	LongDescription  string      `json:"long_description" valid:"-"`
	ImgUrl           string      `json:"img_url" valid:"url"`
	Categories       []*Category `json:"categories" valid:"-"`
}

func (d *RestaurantAndFood) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func NewCategory(c *entity.Category) *Category {
	return &Category{
		Id:   uint64(c.Id),
		Name: c.Name,
		Food: NewFoodInCategoryArr(c.Food),
	}
}

func NewCategoryArray(categories []*entity.Category) []*Category {
	if len(categories) == 0 {
		return make([]*Category, 0)
	}
	cDTO := make([]*Category, len(categories))
	for i, c := range categories {
		cDTO[i] = NewCategory(c)
	}
	return cDTO
}

func NewRestaurantAndFood(r *entity.Restaurant, categories []*entity.Category) *RestaurantAndFood {
	return &RestaurantAndFood{
		Id:               r.Id,
		Name:             r.Name,
		ShortDescription: r.ShortDescription,
		LongDescription:  r.LongDescription,
		ImgUrl:           r.ImgUrl,
		Categories:       NewCategoryArray(categories),
	}
}
