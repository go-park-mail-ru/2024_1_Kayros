package dto

import (
	"2024_1_kayros/internal/entity"
	"github.com/asaskevich/govalidator"
)

type Order struct {
	Id            uint64         `json:"id" valid:"-"`
	UserId        uint64         `json:"user_id" valid:"-"`
	DateOrder     string         `json:"date_order" valid:"-"`
	DateReceiving string         `json:"date_receiving" valid:"-"`
	Status        string         `json:"status" valid:"-"`
	Address       string         `json:"address" valid:"-"`
	ExtraAddress  string         `json:"extra_address" valid:"-"`
	Sum           uint64         `json:"sum" valid:"-"`
	Food          []*FoodInOrder `json:"food" valid:"-"`
}

func (d *Order) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func NewOrder(order *entity.Order) *Order {
	food := order.Food
	foodInOrder := NewFoodArray(food)
	return &Order{
		Id:            order.Id,
		UserId:        order.UserId,
		DateOrder:     order.DateOrder,
		DateReceiving: order.DateReceiving,
		Status:        order.Status,
		Address:       order.Address,
		ExtraAddress:  order.ExtraAddress,
		Sum:           order.Sum,
		Food:          foodInOrder,
	}
}
