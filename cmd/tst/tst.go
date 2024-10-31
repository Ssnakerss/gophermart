package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Ssnakerss/gophermart/internal/db"
	"github.com/Ssnakerss/gophermart/internal/mock"
	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/Ssnakerss/gophermart/internal/types"
)

func main() {
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stor, err := db.New(db.ConString, db.Info)
	if err != nil {
		panic(err)
	}

	orders := []models.Order{
		{Status: types.NEW},
		{Status: types.REGISTERED},
		{Status: types.PROCESSING},
	}
	tcnt := 0
	as := mock.NewMockAccrualService(ctx)

	for {
		pendingOrders := stor.GetOrdersByStatus(ctx, orders)
		if len(pendingOrders) == 0 {
			fmt.Println("no pending orders")
			break // no orders
		}

		wg := &sync.WaitGroup{}
		res := make(chan *models.Order)
		cnt := 0

		go func() {
			for ord := range res {
				time.Sleep(time.Millisecond * 5) // simulate db save
				stor.SaveOrder(ctx, ord)
				fmt.Println(ord) // print accrual
				cnt++
				tcnt++
			}
		}()

		for _, o := range pendingOrders {
			wg.Add(1)
			ord := o
			go func() {
				accrual, err := as.GetAccrual(ord.Number)
				if err != nil {
					wg.Done()
					return
				}
				ord.Accrual = accrual.Accrual
				ord.Status = accrual.Status
				res <- &ord
				wg.Done()
			}()
		}
		wg.Wait() // wait for all goroutines to finish
		close(res)
		fmt.Println("Batch Done", cnt)
		time.Sleep(time.Second * 10)
	}

	fmt.Println("All Done", tcnt)
}
