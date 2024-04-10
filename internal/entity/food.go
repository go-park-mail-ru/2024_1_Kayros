package entity

import "2024_1_kayros/internal/utils/alias"

type Category struct {
	Id   alias.CategoryId
	Name string
	Food []*Food
}

type Food struct {
	Id           uint64
	Name         string
	Description  string
	RestaurantId uint64
	Category     string
	Weight       uint64
	Price        uint64
	ImgUrl       string
}

type FoodInOrder struct {
	Id           uint64
	Name         string
	Weight       uint64
	Price        uint64
	Count        uint64
	ImgUrl       string
	RestaurantId uint64
}
