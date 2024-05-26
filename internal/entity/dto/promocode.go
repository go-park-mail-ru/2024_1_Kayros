package dto

import (
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/constants"
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
		Sale: promo.Sale,
		Date: promo.Date.Format(constants.Timestamptz),
	}
}

func NewPromocodeArray(arr []*entity.Promocode) []*Promocode {
	codeArray := make([]*Promocode, len(arr))
	for i, o := range arr {
		codeArray[i] = NewPromocode(o)
	}
	return codeArray
}
