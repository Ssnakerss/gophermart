package models

import (
	"context"
)

type OrderStorage interface {
	New(ctx context.Context, params ...string) error

	SaveOrder(ctx context.Context, order *Order) error
	GetOrder(ctx context.Context, order *Order) error

	GetOrdersByUser(ctx context.Context, user *User) []Order
	GetOrdersByStatus(ctx context.Context, stat OrderStatus) []Order
}
