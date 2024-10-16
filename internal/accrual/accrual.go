package accrual

import (
	"context"

	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/Ssnakerss/gophermart/internal/types"
)

type AccrualSystem interface {
	GetAccrual(order types.OrderNum) (*models.AccrualResponse, error)
}

type AccrualGetter struct {
	accSystem AccrualSystem
	storage   models.OrderStorage
}

func NewAccrualGetter(accSystem AccrualSystem, storage models.OrderStorage) *AccrualGetter {
	return &AccrualGetter{
		accSystem: accSystem,
		storage:   storage,
	}
}

func (ag *AccrualGetter) GetPendingOrders(ctx context.Context) []models.Order {
	orders := []models.Order{
		{Status: types.NEW},
		{Status: types.REGISTERED},
		{Status: types.PROCESSING},
	}
	return ag.storage.GetOrdersByStatus(ctx, orders)
}

func (ag *AccrualGetter) GetAccrual(ctx context.Context, orderToCheck *models.Order) (*models.AccrualResponse, error) {
	return ag.accSystem.GetAccrual(orderToCheck.Number)
}

func (ag *AccrualGetter) UpdateOrderStatus(ctx context.Context, orderToUpdate *models.Order) error {
	return ag.storage.SaveOrder(ctx, orderToUpdate)
}

func (ag *AccrualGetter) UpdateAccruals(ctx context.Context) error {
	pendingOrders := ag.GetPendingOrders(ctx)
	for _, order := range pendingOrders {
		accrual, err := ag.GetAccrual(ctx, &order)
		if err != nil {
			return err
		}
		if accrual.Status != order.Status {
			order.Status = accrual.Status
			err = ag.UpdateOrderStatus(ctx, &order)
			if err != nil {
				return err
			}
		}
	}
	return nil

}
