package entity

import "database/sql"

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
	Food         []*FoodInOrder
}

type OrderDB struct {
	Id           uint64
	UserId       uint64
	CreatedAt    sql.NullString
	UpdatedAt    sql.NullString
	ReceivedAt   sql.NullString
	Status       string
	Address      string
	ExtraAddress sql.NullString
	Sum          uint64
	Food         []*FoodInOrder
}

func ToOrder(oDB *OrderDB) *Order {
	return &Order{
		Id:           oDB.Id,
		UserId:       oDB.UserId,
		CreatedAt:    String(oDB.CreatedAt),
		UpdatedAt:    String(oDB.UpdatedAt),
		ReceivedAt:   String(oDB.ReceivedAt),
		Status:       oDB.Status,
		Address:      oDB.Address,
		ExtraAddress: String(oDB.ExtraAddress),
		Sum:          oDB.Sum,
		Food:         oDB.Food,
	}
}

func String(element sql.NullString) string {
	if element.Valid {
		return element.String
	}
	return ""
}
