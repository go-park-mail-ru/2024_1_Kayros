package entity

import (
	"database/sql"
)

type Order struct {
	Id           uint64
	UserId       uint64
	CreatedAt    string
	UpdatedAt    string
	ReceivedAt   string
	Status       string
	Address      string
	ExtraAddress string
	Sum          uint64
	RestaurantId uint64
	Food         []*FoodInOrder
}

type OrderDB struct {
	Id           uint64
	UserId       sql.NullInt64
	CreatedAt    string
	UpdatedAt    sql.NullString
	ReceivedAt   sql.NullString
	Status       string
	Address      sql.NullString
	ExtraAddress sql.NullString
	Sum          sql.NullInt64
	Food         []*FoodInOrder
}

func ToOrder(oDB *OrderDB) *Order {
	return &Order{
		Id:           oDB.Id,
		UserId:       Int(oDB.UserId),
		CreatedAt:    oDB.CreatedAt,
		UpdatedAt:    String(oDB.UpdatedAt),
		ReceivedAt:   String(oDB.ReceivedAt),
		Status:       oDB.Status,
		Address:      String(oDB.Address),
		ExtraAddress: String(oDB.ExtraAddress),
		Sum:          Int(oDB.Sum),
		Food:         oDB.Food,
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
