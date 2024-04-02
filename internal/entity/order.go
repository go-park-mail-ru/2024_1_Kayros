package entity

type Order struct {
	Id            uint64
	UserId        uint64
	DateOrder     string
	DateReceiving string
	Status        string
	Address       string
	ExtraAddress  string
	Sum           uint64
	Food          []*FoodInOrder
}
