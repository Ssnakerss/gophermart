package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/Ssnakerss/gophermart/internal/accrual"
	"github.com/Ssnakerss/gophermart/internal/db"
	"github.com/Ssnakerss/gophermart/internal/handlers"
	"github.com/Ssnakerss/gophermart/internal/mock"
	"github.com/Ssnakerss/gophermart/internal/router"
	"golang.org/x/sync/errgroup"
)

func RunWithContext(ctx context.Context, endPoint string) error {
	slog.Info("initialize storage")
	storage, err := db.New(db.ConString, db.Warn)
	if err != nil {
		log.Fatal("db init failed", err)
	}
	slog.Info("migrate data schema")
	storage.Migrate(ctx)
	if err != nil {
		slog.Error(err.Error())
	}

	handlerMaster := handlers.NewMaster(ctx, storage)
	router := router.New(handlerMaster)
	//созджаем сервер
	s := &http.Server{
		Addr:    endPoint,
		Handler: router,
	}

	//создаем обработчик для работы с системой бонусов
	// accrualService := accrual.NewHTTPAccrualsystem("accrual_endpoint")
	accrualService := mock.NewMockAccrualService(ctx)

	ag := accrual.NewAccrualGetter(accrualService, storage)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return ag.Run(ctx, accrual.RunInterval, accrual.BatchSize)
	})

	g.Go(func() error {
		return s.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()
		timeOutCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		return s.Shutdown(timeOutCtx)
	})

	slog.Info("server started")
	return g.Wait()
}
