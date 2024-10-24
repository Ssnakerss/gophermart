package types

import (
	"fmt"
	"strconv"
	"strings"

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
