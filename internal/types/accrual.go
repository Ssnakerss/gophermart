package types

import (
	"fmt"
	"strconv"
	"strings"
)

// --------------------------------------------------------------------------------
// https://golangprojectstructure.com/representing-money-and-currency-go-code/
// тип для хранение бонусов
// использем int
type Accrual int64

func (b *Accrual) MarshalJSON() ([]byte, error) {
	// bb := strconv.FormatFloat(b.Get(), 'f', 2, 64)
	return []byte(fmt.Sprintf(`%f`, b.Get())), nil
}

func (b *Accrual) UnmarshalJSON(bb []byte) error {
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

func (b *Accrual) Set(fb float64) {
	*b = Accrual(uint64(fb * 100))
}
func (b *Accrual) Get() float64 {
	return float64(uint64(*b)) / 100
}
func (b *Accrual) Add(fb float64) {
	b.Set(b.Get() + fb)
}
func (b *Accrual) Sub(fb float64) error {
	f := b.Get() - fb
	if f < 0 {
		return fmt.Errorf("got negative value %f ", f)
	}
	b.Set(f)
	return nil
}
