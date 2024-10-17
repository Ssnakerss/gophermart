package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/theplant/luhn"
)

// номер заказа с дополнительной проверкой luna при присвоении
type OrderNum uint64

func (on *OrderNum) Check() error {
	if !luhn.Valid(int(*on)) || *on == 0 {
		return fmt.Errorf("luna check failed: %d", luhn.CalculateLuhn(int(*on)))
	}
	return nil
}

// сделаем проверку корректности номера заказа в момент присвоения
func (on *OrderNum) Set(num uint64) error {
	*on = OrderNum(num)
	if err := on.Check(); err != nil {
		*on = 0
		return err
	}
	return nil
}

func (on *OrderNum) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"") //убираем кавычки
	if s == "null" {
		return nil
	}
	//пробуем преобразовать в число
	res, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	*on = OrderNum(res)
	return nil
}

func (on *OrderNum) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d"`, *on)), nil
}

func (on *OrderNum) String() string {
	return fmt.Sprintf("%d", *on) //переводим в строку
}

// --------------------------------------------------------------------------------
// кастомный тип для времени чтобы настроит маршаллинг в формат RFC3339
type TimeRFC3339 time.Time

const layout = "2006-01-02T15:04:05Z07:00"

func (t TimeRFC3339) String() string {
	tt := time.Time(t)
	return tt.Format(layout)
}

func (t *TimeRFC3339) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"") //убираем кавычки
	if s == "null" {
		return nil
	}
	tt, err := time.Parse(layout, s)
	*t = TimeRFC3339(tt)
	return err
}

func (t TimeRFC3339) MarshalJSON() ([]byte, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, tt.Format(layout))), nil
}

// --------------------------------------------------------------------------------
// https://golangprojectstructure.com/representing-money-and-currency-go-code/
// тип для хранение бонусов
// использем int
type Bonus int64

func (b *Bonus) MarshalJSON() ([]byte, error) {
	// bb := strconv.FormatFloat(b.Get(), 'f', 2, 64)
	return []byte(fmt.Sprintf(`%f`, b.Get())), nil
}

func (b *Bonus) UnmarshalJSON(bb []byte) error {
	s := strings.Trim(string(bb), "\"") //убираем кавычки
	if s == "null" {
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	b.Set(f)
	return nil
}

func (b *Bonus) Set(fb float64) {
	*b = Bonus(uint64(fb * 100))
}
func (b *Bonus) Get() float64 {
	return float64(uint64(*b)) / 100
}
func (b *Bonus) Add(fb float64) {
	b.Set(b.Get() + fb)
}
func (b *Bonus) Sub(fb float64) error {
	f := b.Get() - fb
	if f < 0 {
		return fmt.Errorf("got negative value %f ", f)
	}
	b.Set(f)
	return nil
}

// --------------------------------------------------------------------------------
// статус заказа с сопоставлением кода ответа
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
	//промежуточные статусы
	NEW        OrderStatus = "NEW"        //заказ загружен но еще не приянт в обработку
	REGISTERED OrderStatus = "REGISTERED" //заказ зарегистрирован но бонус не расчитан
	PROCESSING OrderStatus = "PROCESSING" //заказ в оборботке бонусы расчитыватся

	//финальлные статусы
	INVALID OrderStatus = "INVALID" //система расчета бонусов отказала в расчете

	//такие заказы пишем и в историю по счету
	//Debit-Credit indicator = D
	PROCESSED OrderStatus = "PROCESSED" //заказ обработан, бонусы получены - пишем в дебет счета
)

// статуся для операций с  бонусами
// заказы с этими статусами пишем только в историю по счету
const (
	//это списание бонусов по заказу - пишем в кредит счета
	//Debit-Credit indicator = C
	WITHDRAW OrderStatus = "WITHDRAW"
	//бонусов не хватает списания - для истории
	//Debit-Credit indicator = E
	NOTENOUGHBONUS OrderStatus = "NOTENOUGHBONUS"
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
