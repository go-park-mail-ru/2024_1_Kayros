package dto

import "2024_1_kayros/internal/entity"

type FoodDTO struct {
	Id          uint64 `json:"id" valid:"-"`
	Name        string `json:"name" valid:"-"`
	Description string `json:"description,omitempty" valid:"-"`
	Restaurant  uint64 `json:"restaurant" valid:"-"`
	ImgUrl      string `json:"img_url" valid:"-"`
	Weight      uint64 `json:"weight" valid:"-"`
	Price       uint64 `json:"price" valid:"-"`
}

type FoodInOrder struct {
	Id     uint64 `json:"id" valid:"-"`
	Name   string `json:"name" valid:"-"`
	ImgUrl string `json:"img_url" valid:"-"`
	Weight uint64 `json:"weight" valid:"-"`
	Price  uint64 `json:"price" valid:"-"`
	Count  uint64 `json:"count" valid:"-"`
}

func NewFoodArray(food []*FoodInOrder, orderFood []*entity.FoodInOrder) []*FoodInOrder {
	for i, v := range orderFood {
		food[i].Id = v.Id
		food[i].Name = v.Name
		food[i].ImgUrl = v.ImgUrl
		food[i].Weight = v.Weight
		food[i].Price = v.Price
		food[i].Count = v.Count
	}
	return food
}
