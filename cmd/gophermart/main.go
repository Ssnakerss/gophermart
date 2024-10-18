package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ssnakerss/gophermart/internal/logger"
	"github.com/Ssnakerss/gophermart/internal/server"
)

func main() {
	logger.Setup("DEV")
	slog.Info("server starting")

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

	slog.Warn("server", "status", server.RunWithContext(ctx, ":8080"))
	slog.Info("server stopped")
}
