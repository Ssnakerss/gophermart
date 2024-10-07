package mock

import (
	"context"
	"fmt"
	"math/rand/v2"
)

//int - response code
//string - processing status
//float32 - acrual amount

func YourAcrualIs(order int) (int, string, float32) {
	switch order {
	case 101:
		return 200, "REGISTERED", -1
	case 102:
		return 200, "INVALID", 0
	case 103:
		return 200, "PROCESSING", -1
	case 104:
		return 200, "PROCESSED", 100.500
	case 204:
		return 204, "NOT REGISTERED", -1
	case 429:
		return 429, " No more than N requests per minute allowed", -1
	case 500:
		return 500, "INTERNAL ERROR", -1
	default:
		rnd := randRange(1, 4)
		fmt.Println(rnd)
		if rnd == 2 {
			return 200, "PROCESSED", float32(randRange(1, 1000)) * 1.1
		} else {
			return 200, "PROCESSING", -1
		}
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func StartMockService(ctx context.Context, port int) error {

	return nil
}
