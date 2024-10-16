package mock

import (
	"fmt"
	"math/rand/v2"

	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/Ssnakerss/gophermart/internal/types"
)

//int - response code
//string - processing status
//float32 - acrual amount

type MockAccrualService struct {
}

func (mas *MockAccrualService) GetAcrual(order types.OrderNum) (*models.AccrualResponse, error) {
	ar := models.AccrualResponse{Order: order}
	_, status, accrual := yourAcrualIs(ar.Order)
	ar.Status = types.OrderStatus(status)
	ar.Accrual.Set(accrual)
	return &ar, nil
}

func yourAcrualIs(order types.OrderNum) (int, string, float64) {
	switch order {
	case 1000009:
		return 200, "REGISTERED", -1
	case 1000017:
		return 200, "INVALID", 0
	case 1000025:
		return 200, "PROCESSING", -1
	case 1000033:
		return 200, "PROCESSED", 100.500
	case 1000041:
		return 204, "NOT REGISTERED", -1
	case 1000058:
		return 429, " No more than N requests per minute allowed", -1
	case 1000066:
		return 500, "INTERNAL ERROR", -1
	default:
		rnd := randRange(1, 5)
		fmt.Println(rnd)
		if rnd == 1 {
			return 200, string(types.PROCESSED), float64(randRange(1, 1000)) * 1.1
		}
		if rnd == 2 {
			return 200, string(types.INVALID), -1
		}
		if rnd == 3 {
			return 200, string(types.REGISTERED), -1
		}

		return 200, string(types.PROCESSING), -1
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
