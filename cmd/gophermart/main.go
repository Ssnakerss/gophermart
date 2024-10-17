package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ssnakerss/gophermart/internal/flags"
	"github.com/Ssnakerss/gophermart/internal/logger"
	"github.com/Ssnakerss/gophermart/internal/server"
)

func main() {
	logger.Setup("DEV")
	slog.Info("server starting")

	slog.Info("reading configuration")
	appCfg := flags.NewAppConfig()
	if appCfg == nil {
		slog.Error("failed to parse flags")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		sig := <-exit
		slog.Warn("signal received", "termination", sig)
		slog.Info("stopping server")
		cancel()
	}()

	slog.Warn("server", "status", server.RunWithContext(ctx, appCfg))
	slog.Info("server stopped")
}
