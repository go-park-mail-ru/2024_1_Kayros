package entity

import (
	"database/sql"
)

type ShortOrder struct {
	Id             uint64
	UserId         uint64
	Status         string
	Time           string
	RestaurantId   uint64
	RestaurantName string
	Sum            uint32
}

type Order struct {
	Id             uint64
	UserId         uint64
	CreatedAt      string
	UpdatedAt      string
	ReceivedAt     string
	OrderCreatedAt string
	DeliveredAt    string
	Status         string
	Address        string
	ExtraAddress   string
	Sum            uint64
	NewSum         uint64
	Promocode      string
	RestaurantId   uint64
	RestaurantName string
	Commented      bool
	Food           []*FoodInOrder
}

type OrderDB struct {
	Id             uint64
	UserId         sql.NullInt64
	CreatedAt      string
	UpdatedAt      sql.NullString
	ReceivedAt     sql.NullString
	OrderCreatedAt sql.NullString
	DeliveredAt    sql.NullString
	Status         string
	Address        sql.NullString
	ExtraAddress   sql.NullString
	Sum            sql.NullInt64
	Commented      bool
	Food           []*FoodInOrder
}

func ToOrder(oDB *OrderDB) *Order {
	return &Order{
		Id:             oDB.Id,
		UserId:         Int(oDB.UserId),
		CreatedAt:      oDB.CreatedAt,
		UpdatedAt:      String(oDB.UpdatedAt),
		ReceivedAt:     String(oDB.ReceivedAt),
		OrderCreatedAt: String(oDB.OrderCreatedAt),
		DeliveredAt:    String(oDB.DeliveredAt),
		Status:         oDB.Status,
		Address:        String(oDB.Address),
		ExtraAddress:   String(oDB.ExtraAddress),
		Sum:            Int(oDB.Sum),
		Food:           oDB.Food,
		Commented:      oDB.Commented,
	}
}

func String(element sql.NullString) string {
	if element.Valid {
		return element.String
	}
	return ""
}

func Int(element sql.NullInt64) uint64 {
	if element.Valid {
		return uint64(element.Int64)
	}
	return 0
}
