package models

import "fmt"

//https://golangprojectstructure.com/representing-money-and-currency-go-code/
type Bonus uint64

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

type Account struct {
	Number  string
	Balance Bonus
	History []Transaction
}

type Transaction struct {
	AccountNumber string
	DebitCreadit  string
	OrderNumber   int
	Amount        Bonus
}
