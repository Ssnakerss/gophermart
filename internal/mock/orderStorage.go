package mock

import (
	"context"
	"fmt"

	"github.com/Ssnakerss/gophermart/internal/apperrs"
	"github.com/Ssnakerss/gophermart/internal/models"
)

type mockOrderStorage struct {
	stor map[string]*models.Order
}

func NewMockOrderStorage() *mockOrderStorage {
	m := mockOrderStorage{}
	m.stor = make(map[string]*models.Order)
	return &m
}

func (m *mockOrderStorage) SaveOrder(ctx context.Context, order *models.Order) error {
	_, ok := m.stor[order.Number.String()]
	if ok {
		return fmt.Errorf("Order with number %s already exists", order.Number.String())
	}
	m.stor[order.Number.String()] = order
	return nil
}

func (m *mockOrderStorage) GetOrder(ctx context.Context, order *models.Order) error {

	o, ok := m.stor[order.Number.String()]

	if !ok {
		return apperrs.ErrRecordNotFound
	}
	order.UserID = o.UserID
	order.Accrual = o.Accrual
	order.Status = o.Status
	order.TimeStamp = o.TimeStamp
	order.UpdatedAt = o.UpdatedAt
	return nil
}

func (m *mockOrderStorage) GetAllOrders(ctx context.Context, order *models.Order) []models.Order {
	res := make([]models.Order, len(m.stor))
	for _, v := range m.stor {
		if v.Status == order.Status || v.UserID == order.UserID {
			res = append(res, *v)
		}
	}
	return res
}

func (m *mockOrderStorage) GetOrdersByStatus(ctx context.Context, ordersToFind []models.Order) []models.Order {
	res := make([]models.Order, 0)
	for _, otf := range ordersToFind {
		for _, v := range m.stor {
			if v.Status == otf.Status {
				res = append(res, *v)
			}
		}
	}
	return res
}

func (m *mockOrderStorage) PostTransaction(ctx context.Context, tr *models.Transaction) error {
	return nil
}
