package models

import (
	"context"
)

type OrderStorage interface {
	//вставлят или обновляет заказ в хранилище по номеру
	SaveOrder(ctx context.Context, order *Order) error
	//достает заказ из хранилища по номеру
	GetOrder(ctx context.Context, order *Order) error
	//достает все заказы из хранилища по параметрам
	//в нашем случае это имя пользователя и/или статус
	GetAllOrders(ctx context.Context, order *Order) []Order
}
