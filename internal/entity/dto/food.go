package dto

import (
	"github.com/asaskevich/govalidator"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
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
	Id           uint64 `json:"id" valid:"-"`
	Name         string `json:"name" valid:"-"`
	ImgUrl       string `json:"img_url" valid:"-"`
	Weight       uint64 `json:"weight" valid:"-"`
	Price        uint64 `json:"price" valid:"-"`
	Count        uint64 `json:"count" valid:"-"`
	RestaurantId uint64 `json:"restaurant_id" valid:"-"`
}

func (d *Food) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func (d *FoodInOrder) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func NewFood(f *entity.Food) *Food {
	return &Food{
		Id:          f.Id,
		Name:        f.Name,
		Description: f.Description,
		Restaurant:  f.RestaurantId,
		ImgUrl:      f.ImgUrl,
		Weight:      f.Weight,
		Price:       f.Price,
		Category:    f.Category,
	}
}

func NewFoodInCategoryArr(food []*entity.Food) []*Food {
	if len(food) == 0 {
		return make([]*Food, 0)
	}
	foodList := make([]*Food, len(food))
	for i, f := range food {
		foodList[i] = NewFood(f)
	}
	return foodList
}

func NewFoodInOrder(f *entity.FoodInOrder) *FoodInOrder {
	return &FoodInOrder{
		Id:           f.Id,
		Name:         f.Name,
		ImgUrl:       f.ImgUrl,
		Weight:       f.Weight,
		Price:        f.Price,
		Count:        f.Count,
		RestaurantId: f.RestaurantId,
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
		foodList[i].RestaurantId = food.RestaurantId
	}
	return foodList
}

type FoodCount struct {
	FoodId alias.FoodId `json:"food_id" valid:"positive"`
	Count  uint32       `json:"count" valid:"positive"`
}

func (d *FoodCount) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}
