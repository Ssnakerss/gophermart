package accrual

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/Ssnakerss/gophermart/internal/types"
)

type AccrualService interface {
	GetAccrual(order types.OrderNum) (*models.AccrualResponse, error)
}

type AccrualGetter struct {
	accSystem AccrualService
	storage   models.OrderStorage
}

func NewAccrualGetter(accSystem AccrualService, storage models.OrderStorage) *AccrualGetter {
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

func (ag *AccrualGetter) UpdateAccruals(ctx context.Context, batchSize uint) (uint, error) {
	pendingOrders := ag.GetPendingOrders(ctx)

	slog.Info("accrual getter processing orders", "count", len(pendingOrders))

	for cnt, order := range pendingOrders {
		//говорим что не все заказы были обработаны
		//может увеличить раземер батча?
		if uint(cnt) > batchSize {
			return uint(cnt), models.ErrBatchIncomplete
		}
		accrual, err := ag.GetAccrual(ctx, &order)
		//превышено количество запросов к системе расчета
		//надо уменьшить размер батча
		if errors.Is(err, models.ErrTooManyRequests) {
			return uint(accrual.Accrual), models.ErrTooManyRequests
		}

		//прочие ошибки не обрабатываем?
		if err != nil {
			return batchSize, models.ErrAccSystemProblem
		}

		slog.Info("accrual system response", "order", order.Number, "acc", accrual.Accrual, "status", accrual.Status)

		if accrual.Status != order.Status {
			order.Status = accrual.Status
			err = ag.UpdateOrderStatus(ctx, &order)
			if err != nil {
				return batchSize, err
			}
		}
	}
	return batchSize, nil
}

// начальные размеры батча ставим максмально большим
// т.к. не  знаем емкости системы бонусов
var BatchSize uint = ^uint(0)
var BatchIncreaseStep uint = 100
var RunInterval int = 10 //секунды

// TODO - подумать
// изначально не знаем емкость (мощность) rpm - request per minute системы расчета бонусов
// и время выполнения одного запроса
// если время маленькое - можно обработать всё в одном потомке увелиивая размер батча
// если время значительное - в одном потоке всю емкость системы мы не выберем, надо добавлять воркеров
// чтобы забить систему максимальным количеством запросов в минуту
func (ag *AccrualGetter) Run(ctx context.Context, intervalSec int, batchSize uint) error {
	delayNum := 0
	retryIntervals := []int{0, 5, 10, 20}
	var err error
	var bs uint
	for { //main loop
		slog.Info("starting accrual getter service with interval ", "sec", intervalSec, "batch size", batchSize)
		runPeriod := time.NewTicker(time.Second * time.Duration(intervalSec))
		for { //worker loop
			select {
			case <-runPeriod.C:
				bs, err = ag.UpdateAccruals(ctx, batchSize)
			case <-ctx.Done():
				runPeriod.Stop()
				slog.Info("accrual getter service stopped")
				return ctx.Err() //выход из обрабртчика- завершаем работу
			}
			if err != nil {
				break
			}
		}

		slog.Warn("accrual getter service", "event", err)
		switch {
		case errors.Is(err, models.ErrTooManyRequests):
			//устанавливаем значения батча в соответствии с ответом сервера
			batchSize = bs
			//уменьшаем размер увеличения батча для точной настройки
			if delayNum == 0 {
				BatchIncreaseStep = BatchIncreaseStep / 2
			}
			delayNum++
			if delayNum == len(retryIntervals) {
				delayNum = 0
			}
		case errors.Is(err, models.ErrBatchIncomplete):
			//пробуем увеличить батч
			batchSize = batchSize + BatchIncreaseStep
			delayNum = 0

		case errors.Is(err, models.ErrAccSystemProblem):
			//вынесли в default
			//какие-то проблемы со связью - пробуем позже
		default:
			slog.Error("accrual getter", "error", err.Error())
			delayNum++
			if delayNum == len(retryIntervals) {
				delayNum = 0
			}
		}
		if delayNum > 0 {
			slog.Warn("accrual getter service", "pause operation", retryIntervals[delayNum])
			time.Sleep(time.Duration(retryIntervals[delayNum]) * time.Second)
		}
	}
}
