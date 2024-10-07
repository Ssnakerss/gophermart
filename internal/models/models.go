package models

import "time"

type User struct {
	ID            string
	Jwt           string
	IsAUthorizaed bool
}

type Account struct {
	Number  string
	Balance float32
	History []Transaction
}

type Transaction struct {
	AccountNumber string
	DebitCreadit  string
	OrderNumber   int
	Amount        float32
}

type Order struct {
	Number    int
	UserId    string
	Accrual   float32
	Status    string
	TimeStamp time.Time
}
