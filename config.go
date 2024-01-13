package ft

import (
	"log/slog"
	"runtime/debug"
	"sync/atomic"
)

const (
	instrumentationName = "github.com/amanbolat/ft"
)

var (
	logger         atomic.Value
	logLevel       slog.LevelVar
	tracingEnabled atomic.Bool
	metricsEnabled atomic.Bool

	gitCommit = "unknown_git_commit"
)

func init() {
	logger.Store(slog.Default())
	logLevel.Set(slog.LevelInfo)

	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				gitCommit = setting.Value
			}
		}
	}
}

func DefaultLogger() *slog.Logger {
	return logger.Load().(*slog.Logger)
}

func SetDefaultLogger(l *slog.Logger) {
	logger.Store(l)
}

func SetLogLevel(level slog.Level) {
	logLevel.Set(level)
}

func LogLevel() slog.Level {
	return logLevel.Level()
}

func EnableTracing() {
	tracingEnabled.Store(true)
}

func DisableTracing() {
	tracingEnabled.Store(false)
}

func TracingEnabled() bool {
	return tracingEnabled.Load()
}

func EnableMetrics() {
	metricsEnabled.Store(true)
}

func DisableMetrics() {
	metricsEnabled.Store(false)
}
