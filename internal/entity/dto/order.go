package dto

import (
	"fmt"

	"github.com/asaskevich/govalidator"

	"2024_1_kayros/internal/entity"
)

type Order struct {
	Id           uint64         `json:"id" valid:"-"`
	UserId       uint64         `json:"user_id" valid:"-"`
	CreatedAt    string         `json:"created_at" valid:"-"`
	UpdatedAt    string         `json:"updated_at" valid:"-"`
	ReceivedAt   string         `json:"received_at" valid:"-"`
	Status       string         `json:"status" valid:"-"`
	Address      string         `json:"address" valid:"-"`
	ExtraAddress string         `json:"extra_address" valid:"-"`
	Sum          uint64         `json:"sum" valid:"-"`
	Food         []*FoodInOrder `json:"food" valid:"-"`
}

type FullAddress struct {
	Address      string `json:"address"`
	ExtraAddress string `json:"extra_address"`
}

func (d *Order) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func NewOrder(order *entity.Order) *Order {
	var food []*entity.FoodInOrder
	fmt.Println("dto food len", len(order.Food))
	if len(order.Food) > 0 {
		food = order.Food
	}
	fmt.Println(food)
	foodInOrder := NewFoodArray(food)
	fmt.Println(foodInOrder)
	return &Order{
		Id:           order.Id,
		UserId:       order.UserId,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
		ReceivedAt:   order.ReceivedAt,
		Status:       order.Status,
		Address:      order.Address,
		ExtraAddress: order.ExtraAddress,
		Sum:          order.Sum,
		Food:         foodInOrder,
	}
}

func NewOrders(orderArray []*entity.Order) []*Order {
	orderDTOArray := make([]*Order, len(orderArray))
	for _, order := range orderArray {
		food := order.Food
		foodInOrder := NewFoodArray(food)
		orderDTO := &Order{
			Id:           order.Id,
			UserId:       order.UserId,
			CreatedAt:    order.CreatedAt,
			UpdatedAt:    order.UpdatedAt,
			ReceivedAt:   order.ReceivedAt,
			Status:       order.Status,
			Address:      order.Address,
			ExtraAddress: order.ExtraAddress,
			Sum:          order.Sum,
			Food:         foodInOrder,
		}
		orderDTOArray = append(orderDTOArray, orderDTO)
	}
	return orderDTOArray
}

func NewOrderFromDTO(order *Order) *entity.Order {
	food := order.Food
	foodInOrder := NewFoodArrayFromDTO(food)
	return &entity.Order{
		Id:           order.Id,
		UserId:       order.UserId,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
		ReceivedAt:   order.ReceivedAt,
		Status:       order.Status,
		Address:      order.Address,
		ExtraAddress: order.ExtraAddress,
		Sum:          order.Sum,
		Food:         foodInOrder,
	}
}
