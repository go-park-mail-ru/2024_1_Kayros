package dto

import "2024_1_kayros/internal/entity"

type Order struct {
	Id            uint64         `json:"id"`
	User          uint64         `json:"user"`
	DateOrder     string         `json:"date_order"`
	DateReceiving string         `json:"date_receiving"`
	Status        string         `json:"status"`
	Address       string         `json:"address"`
	ExtraAddress  string         `json:"extra_address"`
	Sum           uint64         `json:"sum"`
	Food          []*FoodInOrder `json:"food"`
}

func OrderToDTO(order *entity.Order, food []*FoodInOrder) *Order {
	return &Order{
		Id:            order.Id,
		User:          order.User,
		DateOrder:     order.DateOrder,
		DateReceiving: order.DateReceiving,
		Status:        order.Status,
		Address:       order.Address,
		ExtraAddress:  order.ExtraAddress,
		Sum:           order.Sum,
		Food:          food,
	}
}
