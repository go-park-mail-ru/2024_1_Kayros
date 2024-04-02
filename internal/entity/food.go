package entity

type Food struct {
	Id           uint64
	Name         string
	Description  string
	RestaurantId uint64
	CategoryId   uint64
	Weight       uint64
	Price        uint64
	ImgUrl       string
}
