package order

import (
	"context"

	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/Ssnakerss/gophermart/internal/user"
)

type OrderProcessor struct {
	storage models.OrderStorage
}

// создаем новый заказ с проверками
// проверка луна - зашита в номер заказа
// так что явно проверяем на дубликаты
func (op *OrderProcessor) New(ctx context.Context, orderNum uint64, user *user.User) *models.Order {
	var o models.Order
	o.UserId = user.Id
	//при создании нового номер если он не проходит проверку луна - выдаст ошибку
	if err := o.Number.Set(orderNum); err != nil {
		o.Status = models.WRONGFORMAT
	} else {
		o.Status = op.CheckIfExist(ctx, &o)
	}
	return &o
}

// проверяем на повторную отправку тем же пользоватем
// и на предмет существования от другого пользователя
func (op *OrderProcessor) CheckIfExist(ctx context.Context, o *models.Order) models.OrderStatus {
	userId := o.UserId
	//если заказа нет в хранилище - вернется ошибка
	//значит заказ новый
	//TODO - разобрать ошибки -  может быть 500 Internal server error
	if err := op.storage.GetOrder(ctx, o); err != nil {
		return models.NEW
	}
	if o.UserId == userId {
		return models.REPEATED
	}
	return models.DUPLICATED
}
