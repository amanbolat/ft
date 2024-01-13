package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amanbolat/ft"
	"github.com/hashicorp/go-metrics"
)

func main() {
	ft.EnableMetrics()
	metricsSink := metrics.NewInmemSink(time.Second*2, time.Minute)

	sig := metrics.DefaultInmemSignal(metricsSink)
	_, _ = metrics.NewGlobal(metrics.DefaultConfig("example"), metricsSink)

	defer sig.Stop()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     ft.LogLevel(),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				source.File = ""
			}

			return a
		},
	}))
	ft.SetLogLevel(slog.Level(1))
	ft.SetDefaultLogger(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	pid := os.Getpid()
	proc, err := os.FindProcess(pid)
	if err != nil {
		os.Exit(1)
	}

	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-ticker.C:
				_ = proc.Signal(syscall.SIGUSR1)
			case <-ctx.Done():
				return
			}
		}
	}()

	ticker := time.NewTicker(time.Millisecond * 300)
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case <-ticker.C:
			_ = Do(ctx)
		}
	}

}

func Do(ctx context.Context) (err error) {
	defer ft.Trace(ctx, "main.Do").WithError(&err).Log()

	err = errors.New("ERROR")

	return
}
