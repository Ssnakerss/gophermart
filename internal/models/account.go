package models

import "github.com/Ssnakerss/gophermart/internal/types"

type Account struct {
	UserID    string            `json:"-" gorm:"primary_key"` //используем UserID вместо номера счета
	Balance   types.Bonus       `json:"current"`              //текущий баланс
	Debit     types.Bonus       `json:"-"`                    //сумма всех поступлений
	Credit    types.Bonus       `json:"withdrawn"`            //сумма всех списаний
	UpdatedAt types.TimeRFC3339 `json:"-"`                    //дата последнего обновления
}

type Transaction struct {
	ID          uint64            `gorm:"primary_key" json:"-"` //идентификатор транзакциив системе
	UserID      string            `json:"-"`                    //используем UserID вместо номера счета
	OrderNumber types.OrderNum    `json:"order"`                //номер заказа по которуму проходила операция
	Bonus       types.Bonus       `json:"sum"`                  //сумма бонусов в операции
	Indicator   string            `json:"-"`                    //Вид операции D-debit +, C-credit - , E-error недостаточно средств для списания
	TimeStamp   types.TimeRFC3339 `json:"processed_at"`
}
