package entity

import "database/sql"

type Promocode struct {
	Id   uint64
	Code string
	Date string
	Sale uint8
	Type string
	Rest uint64
	Sum  uint64
}

//попробовать сделать Type, как enum
// - На первый заказ в сервисе - type=first
// - На первый заказ в ресторане - type=rest
// - После определенной суммы - type=sum
// - Просто единоразовый - type=once

type PromocodeDB struct {
	Id   uint64
	Code string
	Date string
	Sale uint8
	Type string
	Rest sql.NullInt64
	Sum  sql.NullInt64
}

func ToPromocode(promoDB *PromocodeDB) *Promocode {
	return &Promocode{
		Id:   promoDB.Id,
		Code: promoDB.Code,
		Date: promoDB.Date,
		Sale: promoDB.Sale,
		Type: promoDB.Type,
		Rest: Int(promoDB.Rest),
		Sum:  Int(promoDB.Sum),
	}
}