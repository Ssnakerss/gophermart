package order

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/Ssnakerss/gophermart/internal/types"
	"github.com/Ssnakerss/gophermart/internal/user"
)

type OrderProcessor struct {
	storage models.OrderStorage
}

func NewOrderProcessor(st models.OrderStorage) *OrderProcessor {
	var op OrderProcessor
	op.storage = st
	return &op
}

// создаем новый заказ с проверками
// проверка луна - зашита в номер заказа
// так что явно проверяем на дубликаты
func (op *OrderProcessor) NewOrder(ctx context.Context, orderNum string, user *user.User) *models.Order {
	var o models.Order
	o.TimeStamp = types.TimeRFC3339(time.Now())
	ordernum, err := strconv.ParseUint(string(orderNum), 10, 64)
	if err != nil {
		o.Status = types.WRONGFORMAT
		return &o
	}

	o.UserID = user.ID
	//при создании нового номер если он не проходит проверку луна - выдаст ошибку
	if err := o.Number.Set(ordernum); err != nil {
		o.Status = types.WRONGFORMAT
	} else {
		o.Status = op.CheckIfExist(ctx, &o)
	}
	if o.Status == types.NEW {
		op.storage.SaveOrder(ctx, &o)
	}
	return &o
}

func (op *OrderProcessor) AllOrders(ctx context.Context, user *user.User) []models.Order {
	o := models.Order{
		UserID: user.ID,
	}
	return op.storage.GetAllOrders(ctx, &o)
}

// проверяем на повторную отправку тем же пользоватем
// и на предмет существования от другого пользователя
func (op *OrderProcessor) CheckIfExist(ctx context.Context, o *models.Order) types.OrderStatus {
	userID := o.UserID
	//если заказа нет в хранилище - вернется ошибка
	//значит заказ новый
	//TODO - разобрать ошибки -  может быть 500 Internal server error
	if err := op.storage.GetOrder(ctx, o); err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return types.NEW
		}
		return types.ERROR //какая-то другая ошибка, будем считать что Internal Server Error
	}
	if o.UserID == userID {
		return types.REPEATED
	}
	return types.DUPLICATED
}
