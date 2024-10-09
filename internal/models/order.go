package models

import (
	"fmt"
	"time"

	"github.com/theplant/luhn"
)

type Order struct {
	Number    OrderNum    `gorm:"primary_key" json:"number"` //номер заказа
	UserID    string      `gorm:"index"`                     //ид пользователя
	Accrual   Bonus       //размер начисленного бонуса
	Status    OrderStatus `gorm:"index"` //статус заказа
	TimeStamp time.Time   //время поступления заказа
}
type OrderNum uint64

// сделаем проверку корректности номера заказа в момент присвоения
func (on *OrderNum) Set(num uint64) error {
	if !luhn.Valid(int(num)) {
		return fmt.Errorf("luna check failed: %d", luhn.CalculateLuhn(int(num)))
	}
	*on = OrderNum(num)
	return nil
}

type OrderStatus string

// статусы по итогам обработки нового заказа
const (
	CHECKING OrderStatus = "CHECKING" //проверка в процессе

	REPEATED    OrderStatus = "REPEATED"    //заказ уже был загружен этим пользователем
	DUPLICATED  OrderStatus = "DUPLICATED"  //заказ уже был загружен другим пользователем
	WRONGFORMAT OrderStatus = "WRONGFORMAT" //неправильный формат номера заказа

	ERROR OrderStatus = "ERROR"
)

// Статусы принятых заказаов
const (
	NEW        OrderStatus = "NEW"        //заказ загружен но еще не приянт в обработку
	PROCESSING OrderStatus = "PROCESSING" //заказ в оборботке бонусы расчитыватся
	INVALID    OrderStatus = "INVALID"    //система расчета бонусов отказала в расчете
	PROCESSED  OrderStatus = "PROCESSED"  //заказ обработан, бонусы получены
)

func (os OrderStatus) ResponseCode() int {
	switch os {
	case REPEATED:
		return 200
	case NEW:
		return 202
	case DUPLICATED:
		return 409
	case WRONGFORMAT:
		return 422
	case ERROR:
		return 500
	default:
		return 0
	}
}
