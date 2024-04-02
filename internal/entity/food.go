package entity

type Food struct {
	Id          uint64
	Name        string
	Description string
	Restaurant  uint64
	ImgUrl      string
	Weight      uint64
	Price       uint64
}

type FoodInOrder struct {
	Id     uint64
	Name   string
	ImgUrl string
	Weight uint64
	Price  uint64
	Count  uint64
}
