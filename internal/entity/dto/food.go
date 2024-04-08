package dto

import (
	"fmt"

	"github.com/asaskevich/govalidator"

	"2024_1_kayros/internal/entity"
)

type Food struct {
	Id          uint64 `json:"id" valid:"-"`
	Name        string `json:"name" valid:"-"`
	Description string `json:"description" valid:"-"`
	Restaurant  uint64 `json:"restaurant" valid:"-"`
	ImgUrl      string `json:"img_url" valid:"-"`
	Weight      uint64 `json:"weight" valid:"-"`
	Price       uint64 `json:"price" valid:"-"`
	Category    string `json:"category" valid:"-"`
}

type FoodInOrder struct {
	Id     uint64 `json:"id" valid:"-"`
	Name   string `json:"name" valid:"-"`
	ImgUrl string `json:"img_url" valid:"-"`
	Weight uint64 `json:"weight" valid:"-"`
	Price  uint64 `json:"price" valid:"-"`
	Count  uint64 `json:"count" valid:"-"`
}

func (d *Food) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func (d *FoodInOrder) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func NewFoodInOrder(f *entity.FoodInOrder) *FoodInOrder {
	return &FoodInOrder{
		Id:     f.Id,
		Name:   f.Name,
		ImgUrl: f.ImgUrl,
		Weight: f.Weight,
		Price:  f.Price,
		Count:  f.Count,
	}
}

func NewFoodArray(orderFood []*entity.FoodInOrder) []*FoodInOrder {
	if len(orderFood) == 0 {
		return make([]*FoodInOrder, 0)
	}
	foodList := make([]*FoodInOrder, len(orderFood))
	for i, food := range orderFood {
		foodList[i] = NewFoodInOrder(food)
	}
	fmt.Println(foodList[0], len(foodList))
	return foodList
}

func NewFoodArrayFromDTO(orderFood []*FoodInOrder) []*entity.FoodInOrder {
	foodList := make([]*entity.FoodInOrder, len(orderFood))
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
