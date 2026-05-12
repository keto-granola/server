package main

import (
	"log/slog"
	"os"
)

func main() {
	if err := run(); err != nil {
		slog.Error("run", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("shutting down gracefully...")
}

func run() error {
	return nil
}
