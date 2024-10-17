package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Ssnakerss/gophermart/internal/accrual"
	"github.com/Ssnakerss/gophermart/internal/db"
	"github.com/Ssnakerss/gophermart/internal/flags"
	"github.com/Ssnakerss/gophermart/internal/handlers"
	"github.com/Ssnakerss/gophermart/internal/mock"
	"github.com/Ssnakerss/gophermart/internal/router"
	"golang.org/x/sync/errgroup"
)

func RunWithContext(ctx context.Context, cfg *flags.AppConfig) error {
	slog.Info("startign server", "config", cfg)

	slog.Info("initialize storage")
	var dbLogLevel db.LogLevel
	switch cfg.ENV {
	case "DEV": //для разработки
		dbLogLevel = db.Info
	case "PROD": //для продакшена
		dbLogLevel = db.Warn
	}
	slog.Info("initialize database", "level", dbLogLevel)
	storage, err := db.New(cfg.DatabaseURI, dbLogLevel)
	if err != nil {
		slog.Error("failed to initialize database", "error", err.Error())
		os.Exit(1)
	}
	slog.Info("migrate data schema")
	storage.Migrate(ctx)
	if err != nil {
		slog.Error("failed to migrate data schema", "error", err.Error())
		os.Exit(1)
	}

	handlerMaster := handlers.NewMaster(ctx, storage)
	router := router.New(handlerMaster)
	//созджаем сервер
	s := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}

	//создаем обработчик для работы с системой бонусов
	var accrualService accrual.AccrualService
	switch cfg.ENV {
	case "DEV": //для разработки
		accrualService = mock.NewMockAccrualService(ctx)
	case "PROD": //для продакшена
		accrualService = accrual.NewHTTPAccrualsystem(cfg.AccrualSystemAddress)
	}
	//
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
