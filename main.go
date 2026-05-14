package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/keto-granola/server/internal/config"
	productadmin "github.com/keto-granola/server/internal/product/admin"
	"github.com/keto-granola/server/internal/server"
	"github.com/keto-granola/server/internal/store"
)

func main() {
	if err := run(); err != nil {
		slog.Error("run", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.ParseEnv()
	if err != nil {
		return fmt.Errorf("parse env vars %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))
	slog.SetDefault(logger)

	dataStore, err := store.New(ctx)
	if err != nil {
		return fmt.Errorf("create store %v", err)
	}
	defer dataStore.Close()

	instance := server.New(ctx, cfg.Environment, cfg.ClientURL, composeHandlers(dataStore))

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- instance.Start(cfg.Port)
	}()

	select {
	case <-ctx.Done():
		slog.Info("context cancelled")
	case err = <-serverErr:
	}

	if shutdownErr := instance.Stop(); shutdownErr != nil {
		slog.Error("server shutdown", slog.Any("error", shutdownErr))
	}

	return err
}

func composeHandlers(db *store.Store) *server.Handlers {
	return &server.Handlers{
		ProductAdmin: productadmin.NewHandler(productadmin.NewService(db.ProductStore())),
	}
}
