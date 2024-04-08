package entity

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
	Id     uint64
	Name   string
	Weight uint64
	Price  uint64
	Count  uint64
	ImgUrl string
}
