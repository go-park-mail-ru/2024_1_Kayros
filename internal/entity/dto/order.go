package dto

import (
	"github.com/asaskevich/govalidator"

	"2024_1_kayros/internal/entity"
)

type Order struct {
	Id             uint64         `json:"id" valid:"-"`
	UserId         uint64         `json:"user_id" valid:"-"`
	CreatedAt      string         `json:"-" valid:"-"`
	UpdatedAt      string         `json:"-" valid:"-"`
	ReceivedAt     string         `json:"-" valid:"-"`
	OrderCreatedAt string         `json:"created_at,omitempty" valid:"-"`
	DeliveredAt    string         `json:"delivered_at,omitempty" valid:"-"`
	Status         string         `json:"status" valid:"-"`
	Address        string         `json:"address" valid:"-"`
	ExtraAddress   string         `json:"extra_address" valid:"-"`
	Sum            uint64         `json:"sum" valid:"-"`
	NewSum         uint64         `json:"new_sum,omitempty"`
	Promocode      string         `json:"promocode,omitempty"`
	RestaurantId   uint64         `json:"restaurant_id"`
	RestaurantName string         `json:"restaurant_name"`
	Commented      bool           `json:"commented"`
	Error          string         `json:"error,omitempty"`
	Food           []*FoodInOrder `json:"food" valid:"-"`
}

type ShortOrder struct {
	Id             uint64 `json:"id" valid:"-"`
	UserId         uint64 `json:"user_id" valid:"-"`
	Status         string `json:"status" valid:"-"`
	Time           string `json:"time" valid:"-"`
	RestaurantId   uint64 `json:"restaurant_id,omitempty" valid:"-"`
	RestaurantName string `json:"restaurant_name" valid:"-"`
	Sum            uint32 `json:"sum,omitempty" valid:"-"`
}

type ShortOrderArray struct {
	Payload []*ShortOrder `json:"payload" valid:"-"`
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
		Id:     order.Id,
		UserId: order.UserId,
		//CreatedAt:      order.CreatedAt,
		//UpdatedAt:      order.UpdatedAt,
		//ReceivedAt:     order.ReceivedAt,
		OrderCreatedAt: order.OrderCreatedAt,
		DeliveredAt:    order.DeliveredAt,
		Status:         order.Status,
		Address:        order.Address,
		ExtraAddress:   order.ExtraAddress,
		Sum:            order.Sum,
		NewSum:         order.NewSum,
		Promocode:      order.Promocode,
		RestaurantId:   order.RestaurantId,
		RestaurantName: order.RestaurantName,
		Commented:      order.Commented,
		Error:          order.Error,
		Food:           foodInOrder,
	}
}

func NewOrders(orderArray []*entity.Order) []*Order {
	orderDTOArray := make([]*Order, len(orderArray))
	for i, order := range orderArray {
		food := order.Food
		foodInOrder := NewFoodArray(food)
		orderDTO := &Order{
			Id:     order.Id,
			UserId: order.UserId,
			//CreatedAt:    order.CreatedAt,
			//UpdatedAt:    order.UpdatedAt,
			//ReceivedAt:   order.ReceivedAt,
			OrderCreatedAt: order.OrderCreatedAt,
			DeliveredAt:    order.DeliveredAt,
			Status:         order.Status,
			Address:        order.Address,
			ExtraAddress:   order.ExtraAddress,
			Sum:            order.Sum,
			RestaurantId:   order.RestaurantId,
			Commented:      order.Commented,
			Food:           foodInOrder,
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
		Commented:    order.Commented,
		Food:         foodInOrder,
	}
}

func NewShortOrder(order *entity.ShortOrder) *ShortOrder {
	return &ShortOrder{
		Id:             order.Id,
		UserId:         order.UserId,
		Status:         order.Status,
		Time:           order.Time,
		RestaurantId:   order.RestaurantId,
		RestaurantName: order.RestaurantName,
		Sum:            order.Sum,
	}
}

func NewShortOrderArray(arr []*entity.ShortOrder) []*ShortOrder {
	orderArray := make([]*ShortOrder, len(arr))
	for i, o := range arr {
		orderArray[i] = NewShortOrder(o)
	}
	return orderArray
}
