package mock

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync/atomic"
	"time"

	"github.com/Ssnakerss/gophermart/internal/apperrs"
	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/Ssnakerss/gophermart/internal/types"
)

// int - response code
// string - processing status
// float32 - acrual amount
var requestCounter int32

type MockAccrualService struct {
}

func NewMockAccrualService(ctx context.Context) *MockAccrualService {
	counterResetTimer := time.NewTicker(time.Second * 60)
	fmt.Println("MOCK STARTED!!!!!!!!!!!!!!!!!!")
	go func() {
		for {
			select {
			case <-counterResetTimer.C:
				atomic.StoreInt32(&requestCounter, 0)
				fmt.Println("MOCK >>> reset counter")
			case <-ctx.Done():
				counterResetTimer.Stop()
				fmt.Println("MOCK >>> stop reset goroutine")
				return
			}
		}
	}()

	// counterReprtTimer := time.NewTicker(time.Second)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-counterReprtTimer.C:
	// 			fmt.Println("MOCK >>>  counter value: ", atomic.LoadInt32(&requestCounter))
	// 		case <-ctx.Done():
	// 			counterReprtTimer.Stop()
	// 			fmt.Println("MOCK >>> stop report goroutine")
	// 			return
	// 		}
	// 	}
	// }()

	return &MockAccrualService{}
}

func (mas *MockAccrualService) GetAccrual(order types.OrderNum) (*models.AccrualResponse, error) {
	ar := models.AccrualResponse{Order: order}

	maxRequestCnt := int32(randRange(100, 150))
	atomic.AddInt32(&requestCounter, 1)
	if atomic.LoadInt32(&requestCounter) > maxRequestCnt {
		ar.Order = 0
		ar.Status = "E"
		ar.Accrual = types.Accrual(maxRequestCnt)
		return &ar, apperrs.ErrTooManyRequests
	}

	_, status, accrual := yourAcrualIs(ar.Order)
	ar.Status = types.OrderStatus(status)
	ar.Accrual.Set(accrual)
	return &ar, nil
}

func yourAcrualIs(order types.OrderNum) (int, string, float64) {
	time.Sleep(time.Duration(rand.IntN(50)) * time.Millisecond) //типа работаемс задержкой

	switch order {
	case 1000009:
		return 200, "REGISTERED", 0
	case 1000017:
		return 200, "INVALID", 0
	case 1000025:
		return 200, "PROCESSING", 0
	case 1000033:
		return 200, "PROCESSED", 100.500
	case 1000041:
		return 204, "NOT REGISTERED", 0
	case 1000058:
		return 429, " No more than N requests per minute allowed", -1
	case 1000066:
		return 500, "INTERNAL ERROR", 0
	default:
		rnd := randRange(1, 5)

		if rnd == 1 {
			return 200, string(types.PROCESSED), float64(randRange(1, 1000)) * 1.1
		}
		if rnd == 2 {
			return 200, string(types.INVALID), 0
		}
		if rnd == 3 {
			return 200, string(types.REGISTERED), 0
		}

		return 200, string(types.PROCESSING), 0
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
