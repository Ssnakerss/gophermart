package models

import "github.com/Ssnakerss/gophermart/internal/types"

type Order struct {
	Number    types.OrderNum    `gorm:"primary_key" json:"number"` //номер заказа
	UserID    string            `gorm:"index"`                     //ид пользователя
	Accrual   types.Bonus       `json:"accrual"`                   //размер начисленного бонуса
	Status    types.OrderStatus `gorm:"index" json:"status"`       //статус заказа
	TimeStamp types.TimeRFC3339 `json:"uploaded_at"`               //время поступления заказа
}
