package order

import (
	"context"
	"testing"
	"time"

	"github.com/Ssnakerss/gophermart/internal/mock"
	"github.com/Ssnakerss/gophermart/internal/types"

	"github.com/Ssnakerss/gophermart/internal/models"
)

func TestOrderProcessor_NewOrder(t *testing.T) {
	mockStorage := mock.NewMockOrderStorage()
	op := NewOrderProcessor(mockStorage)
	testUser := &models.User{ID: "test", Hash: "test", UpdatedAt: time.Now()}
	anotherTestUser := &models.User{ID: "anothertest", Hash: "anothertest", UpdatedAt: time.Now()}

	type args struct {
		ctx      context.Context
		orderNum string
		user     *models.User
	}
	tests := []struct {
		name string
		op   *OrderProcessor
		args args
		want types.OrderStatus
	}{
		{
			name: "wrong format",
			op:   op,
			args: args{
				ctx:      context.Background(),
				orderNum: "111",
				user:     testUser,
			},
			want: types.WRONGFORMAT,
		},
		{
			name: "good order",
			op:   op,
			args: args{
				ctx:      context.Background(),
				orderNum: "1000256",
				user:     testUser,
			},
			want: types.NEW,
		},
		{
			name: "repeated order",
			op:   op,
			args: args{
				ctx:      context.Background(),
				orderNum: "1000256",
				user:     testUser,
			},
			want: types.REPEATED,
		},
		{
			name: "duplicated order",
			op:   op,
			args: args{
				ctx:      context.Background(),
				orderNum: "1000256",
				user:     anotherTestUser,
			},
			want: types.DUPLICATED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newOrder := tt.op.NewOrder(tt.args.ctx, tt.args.orderNum, tt.args.user)
			if newOrder.Status != tt.want {
				t.Errorf("OrderProcessor.NewOrder() = %v, want %v", newOrder.Status, tt.want)
			}

		})
	}
}
