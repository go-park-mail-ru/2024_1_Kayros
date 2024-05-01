package dto

import (
	"github.com/asaskevich/govalidator"

	"2024_1_kayros/internal/entity"
)

type Order struct {
	Id           uint64         `json:"id" valid:"-"`
	UserId       uint64         `json:"user_id" valid:"-"`
	CreatedAt    string         `json:"created_at" valid:"-"`
	UpdatedAt    string         `json:"-" valid:"-"`
	ReceivedAt   string         `json:"received_at,omitempty" valid:"-"`
	Status       string         `json:"status" valid:"-"`
	Address      string         `json:"address" valid:"-"`
	ExtraAddress string         `json:"extra_address" valid:"-"`
	Sum          uint64         `json:"sum" valid:"-"`
	RestaurantId uint64         `json:"restaurant_id"`
	Food         []*FoodInOrder `json:"food" valid:"-"`
}

type FullAddress struct {
	Address      string `json:"address" valid:"user_address_domain"`
	ExtraAddress string `json:"extra_address" valid:"user_extra_address_domain"`
}

func (d *FullAddress) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func (d *Order) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

func NewOrder(order *entity.Order) *Order {
	food := []*entity.FoodInOrder{}
	if len(order.Food) > 0 {
		food = order.Food
	}
	foodInOrder := NewFoodArray(food)
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
		RestaurantId: order.RestaurantId,
		Food:         foodInOrder,
	}
}

func NewOrders(orderArray []*entity.Order) []*Order {
	orderDTOArray := make([]*Order, len(orderArray))
	for i, order := range orderArray {
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
			RestaurantId: order.RestaurantId,
			Food:         foodInOrder,
		}
		orderDTOArray[i] = orderDTO
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
		RestaurantId: order.RestaurantId,
		Food:         foodInOrder,
	}
}
