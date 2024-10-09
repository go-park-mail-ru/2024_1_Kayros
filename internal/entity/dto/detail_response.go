package dto

import "2024_1_kayros/internal/utils/alias"

type ResponseDetail struct {
	Detail string `json:"detail" valid:"-"`
}

type ResponseUrlPay struct {
	Url string `json:"url" valid:"-"`
}

type PayedOrderInfo struct {
	Id     alias.OrderId `json:"id"`
	Status string        `json:"status"`
}
