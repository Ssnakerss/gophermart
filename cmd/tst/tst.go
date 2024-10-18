package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Ssnakerss/gophermart/internal/accrual"
	"github.com/Ssnakerss/gophermart/internal/db"
	"github.com/Ssnakerss/gophermart/internal/logger"
	"github.com/Ssnakerss/gophermart/internal/mock"
	"github.com/Ssnakerss/gophermart/internal/models"
)

func main() {
	logger.Setup("DEV")
	slog.Info("Hello", "module", "tst")

	d, _ := db.New(db.ConString, db.Info)
	ag := accrual.NewAccrualGetter(nil, d)
	po := ag.GetPendingOrders(context.Background())
	fmt.Println(po)

	m := mock.MockAccrualService{}

	order := models.Order{}
	order.Number.Set(1000074)

	for i := 0; i < 50; i++ {
		fmt.Println(m.GetAcrual(order.Number))
	}

}
