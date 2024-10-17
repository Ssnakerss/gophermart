package models

import "github.com/Ssnakerss/gophermart/internal/types"

type AccrualResponse struct {
	Order   types.OrderNum    `json:"order"`
	Status  types.OrderStatus `json:"status"`
	Accrual types.Bonus       `json:"accrual"`
}
