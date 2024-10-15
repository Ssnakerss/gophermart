package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ssnakerss/gophermart/internal/server"
)

func main() {
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

	slog.Info("server starting")
	slog.Warn("server", "status", server.RunWithContext(ctx, ":8080"))
	slog.Info("server stopped")
}
