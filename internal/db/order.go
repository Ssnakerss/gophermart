package db

import (
	"context"
	"errors"

	"github.com/Ssnakerss/gophermart/internal/apperrs"
	"github.com/Ssnakerss/gophermart/internal/models"
	"gorm.io/gorm"
)

func (db *GormDB) SaveOrder(ctx context.Context, order *models.Order) error {
	return db.DB.WithContext(ctx).Save(order).Error
}

func (db *GormDB) GetOrder(ctx context.Context, order *models.Order) error {
	err := db.DB.WithContext(ctx).First(order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrs.ErrRecordNotFound
	}
	return err
}

func (db *GormDB) GetAllOrders(ctx context.Context, order *models.Order) []models.Order {
	var orders []models.Order
	db.DB.WithContext(ctx).Order("time_stamp DESC").Find(&orders, order)
	return orders
}

func (db *GormDB) GetOrdersByStatus(ctx context.Context, ordersToFind []models.Order) []models.Order {
	var result []models.Order
	var orders []models.Order
	for _, o := range ordersToFind {
		db.DB.WithContext(ctx).Order("time_stamp ASC").Find(&orders, o)
		result = append(result, orders...)
	}
	return result
}
