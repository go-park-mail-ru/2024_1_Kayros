package dto

import (
	"2024_1_kayros/internal/entity"
	"github.com/asaskevich/govalidator"
)

type FoodInOrder struct {
	Id     uint64 `json:"id" valid:"-"`
	Name   string `json:"name" valid:"-"`
	ImgUrl string `json:"img_url" valid:"-"`
	Weight uint64 `json:"weight" valid:"-"`
	Price  uint64 `json:"price" valid:"-"`
	Count  uint64 `json:"count" valid:"-"`
}

func (d *FoodInOrder) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func NewFoodArray(orderFood []*entity.FoodInOrder) []*FoodInOrder {
	foodList := make([]*FoodInOrder, 0, len(orderFood)+1)
	for i, food := range orderFood {
		foodList[i].Id = food.Id
		foodList[i].Name = food.Name
		foodList[i].ImgUrl = food.ImgUrl
		foodList[i].Weight = food.Weight
		foodList[i].Price = food.Price
		foodList[i].Count = food.Count
	}
	return foodList
}
