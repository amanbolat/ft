package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/amanbolat/ft"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       ft.LogLevel(),
		ReplaceAttr: nil,
	}))
	ft.SetLogLevel(slog.Level(1))
	ft.SetDefaultLogger(logger)
	ctx := context.Background()

	Do(ctx)
}

func Do(ctx context.Context) {
	defer ft.Trace(ctx, "main.Do").Log()
}
