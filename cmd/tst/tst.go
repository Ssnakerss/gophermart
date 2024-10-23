package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Ssnakerss/gophermart/internal/logger"
	"github.com/Ssnakerss/gophermart/internal/mock"
	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/Ssnakerss/gophermart/internal/proc/order"
)

// func main() {
// 	logger.Setup("DEV")
// 	// slog.Info("Hello", "module", "tst")

// 	// d, _ := db.New(db.ConString, db.Info)
// 	// ag := accrual.NewAccrualGetter(nil, d)
// 	// po := ag.GetPendingOrders(context.Background())
// 	// fmt.Println(po)

// 	m := mock.NewMockOrderStorage()
// 	op := order.NewOrderProcessor(m)
// 	testUser := &models.User{ID: "test", Hash: "test", UpdatedAt: time.Now()}
// 	anotherTestUser := &models.User{ID: "anothertest", Hash: "anothertest", UpdatedAt: time.Now()}

// 	// o := &models.Order{
// 	// 	Number: 1000256,
// 	// 	UserID: "testy",
// 	// 	Status: types.NEW,
// 	// }

// 	// oo := &models.Order{
// 	// 	Number: 1000256,
// 	// }

// 	newOrder := op.NewOrder(context.Background(), "1000256", testUser)
// 	fmt.Println("-----------------------")
// 	newOrder = op.NewOrder(context.Background(), "1000256", anotherTestUser)

// 	fmt.Println(newOrder)
}
