package ft

import (
	"log/slog"
	"os"

	"github.com/jonboulle/clockwork"
	"github.com/samber/lo"
	"go.uber.org/atomic"
)

var (
	globalLogger             = atomic.NewPointer(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	globalTracingEnabled     = atomic.NewBool(false)
	globalMetricsEnabled     = atomic.NewBool(false)
	globalAppendOtelAttrs    = atomic.NewBool(false)
	globalDurationMetricUnit = atomic.NewString(DurationMetricUnitMillisecond)
	globalClock              = atomic.NewPointer[clockwork.Clock](lo.ToPtr(clockwork.NewRealClock()))

	globalLogLevelEndOnSuccess slog.LevelVar
	globalLogLevelEndOnFailure slog.LevelVar
)

func init() {
	globalLogLevelEndOnSuccess.Set(slog.LevelInfo)
	globalLogLevelEndOnFailure.Set(slog.LevelError)
}

func SetDurationMetricUnit(unit string) {
	switch unit {
	case DurationMetricUnitMillisecond, DurationMetricUnitSecond:
		globalDurationMetricUnit.Store(unit)
	default:
		globalDurationMetricUnit.Store(DurationMetricUnitMillisecond)
	}
}

func SetDefaultLogger(l *slog.Logger) {
	if l == nil {
		return
	}

	globalLogger.Store(l)
}

func SetLogLevelOnFailure(level slog.Level) {
	globalLogLevelEndOnFailure.Set(level)
}

func SetLogLevelOnSuccess(level slog.Level) {
	globalLogLevelEndOnSuccess.Set(level)
}

func SetTracingEnabled(v bool) {
	globalTracingEnabled.Store(v)
}

func SetMetricsEnabled(v bool) {
	globalMetricsEnabled.Store(v)
}

func SetClock(c clockwork.Clock) {
	globalClock.Store(&c)
}

func SetAppendOtelAttrs(v bool) {
	globalAppendOtelAttrs.Store(v)
}
