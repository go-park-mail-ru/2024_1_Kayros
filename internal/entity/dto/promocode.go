package dto

import (
	"2024_1_kayros/internal/entity"
)

type Promocode struct {
	Id   uint64 `json:"id"`
	Code string `json:"code"`
	Date string `json:"date"`
	Sale uint8  `json:"sale"`
}

func NewPromocode(promo *entity.Promocode) *Promocode {
	return &Promocode{
		Id:   promo.Id,
		Code: promo.Code,
		Date: promo.Date,
		Sale: promo.Sale,
	}
}
