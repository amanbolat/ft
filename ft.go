package ft

import (
	"context"
	"fmt"
	"runtime"

	"github.com/puzpuzpuz/xsync/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"

	"log/slog"
	"time"
)

const (
	instrumentationName           = "github.com/amanbolat/ft"
	DurationMetricUnitSecond      = "ms"
	DurationMetricUnitMillisecond = "s"
)

var int64Counters = xsync.NewMapOf[string, metric.Int64Counter]()
var durationHistograms = xsync.NewMapOf[string, metric.Float64Histogram]()

type SpanConfig struct {
	err             *error
	additionalAttrs []slog.Attr
}

type Option func(cfg *SpanConfig)

func WithErr(err *error) Option {
	return func(cfg *SpanConfig) {
		cfg.err = err
	}
}

func WithAttrs(attrs ...slog.Attr) Option {
	return func(cfg *SpanConfig) {
		cfg.additionalAttrs = append(cfg.additionalAttrs, attrs...)
	}
}

type Span interface {
	End()
}

type span struct {
	ctx             context.Context
	start           time.Time
	action          string
	traceSpan       trace.Span
	err             *error
	additionalAttrs []slog.Attr
}

func Start(ctx context.Context, action string, opts ...Option) (context.Context, Span) {
	now := (*globalClock.Load()).Now()

	cfg := &SpanConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	if ctx == nil {
		ctx = context.Background()
	}
	var otelSpan trace.Span

	if globalTracingEnabled.Load() {
		ctx, otelSpan = otel.Tracer(
			instrumentationName,
			trace.WithSchemaURL(semconv.SchemaURL),
		).Start(
			ctx,
			action,
			trace.WithSpanKind(trace.SpanKindInternal),
			trace.WithAttributes(
				attribute.String("action", action),
			),
			trace.WithTimestamp(now),
		)
	}

	if otelSpan != nil && otelSpan.IsRecording() && globalAppendOtelAttrs.Load() && len(cfg.additionalAttrs) > 0 {
		otelAttrs := make([]attribute.KeyValue, 0, len(cfg.additionalAttrs))
		for _, attr := range cfg.additionalAttrs {
			otelAttrs = append(otelAttrs, mapSlogAttrToOtel(attr))
		}
		otelSpan.SetAttributes(otelAttrs...)
	}

	if globalMetricsEnabled.Load() {
		metricName := action + "_counter"
		counter, ok := int64Counters.Load(metricName)

		if !ok {
			var err error
			counter, err = otel.GetMeterProvider().Meter(instrumentationName).Int64Counter(metricName)
			if err == nil {
				int64Counters.Store(metricName, counter)
				ok = true
			}
		}

		if ok {
			counter.Add(ctx, 1)
		}
	}

	attrs := make([]slog.Attr, 1+len(cfg.additionalAttrs))
	attrs = append(attrs, slog.String("action", action))
	attrs = append(attrs, cfg.additionalAttrs...)

	log(ctx, "action started", globalLogLevelEndOnSuccess.Level(), now, attrs...)

	return ctx, span{
		ctx:             ctx,
		start:           now,
		action:          action,
		traceSpan:       otelSpan,
		err:             cfg.err,
		additionalAttrs: cfg.additionalAttrs,
	}
}

func (s span) End() {
	now := (*globalClock.Load()).Now()
	duration := now.Sub(s.start)
	level := globalLogLevelEndOnSuccess.Level()

	durationAttrKey := "duration_ms"
	durationAttrVal := durationToMillisecond(duration)

	attrs := make([]slog.Attr, 2+len(s.additionalAttrs))
	attrs = append(attrs, slog.String("action", s.action), slog.Float64(durationAttrKey, durationAttrVal))
	attrs = append(attrs, s.additionalAttrs...)

	if s.err != nil && *s.err != nil {
		level = globalLogLevelEndOnFailure.Level()
		attrs = append(attrs, slog.Any("error", *s.err))

		if s.traceSpan != nil {
			s.traceSpan.RecordError(*s.err, trace.WithStackTrace(true))
			s.traceSpan.SetStatus(codes.Error, (*s.err).Error())
		}
	}

	if globalMetricsEnabled.Load() {
		durationMetricSuffix := "_duration_milliseconds"
		durationMetricUnit := globalDurationMetricUnit.Load()

		if durationMetricUnit == DurationMetricUnitSecond {
			durationAttrKey = "duration_s"
			durationAttrVal = durationToSecond(duration)
			durationMetricSuffix = "_duration_seconds"
		}

		metricName := s.action + durationMetricSuffix
		histogram, ok := durationHistograms.Load(metricName)

		if !ok {
			var err error
			histogram, err = otel.GetMeterProvider().
				Meter(instrumentationName).
				Float64Histogram(
					metricName,
					metric.WithUnit(durationMetricUnit),
					metric.WithDescription(fmt.Sprintf("[%s] action duration", s.action)),
				)
			if err == nil {
				durationHistograms.Store(metricName, histogram)
				ok = true
			}
		}

		if ok {
			histogram.Record(s.ctx, duration.Seconds())
		}
	}

	log(s.ctx, "action ended", level, now, attrs...)

	if s.traceSpan != nil {
		s.traceSpan.End(trace.WithTimestamp(now))
	}
}

func log(ctx context.Context, msg string, level slog.Level, now time.Time, attrs ...slog.Attr) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(now, level, msg, pcs[0])
	r.AddAttrs(attrs...)
	_ = globalLogger.Load().Handler().Handle(ctx, r)
}

func mapSlogAttrToOtel(v slog.Attr) attribute.KeyValue {
	key := v.Key
	value := v.Value

	switch value.Kind() {
	case slog.KindBool:
		return attribute.Bool(key, value.Bool())
	case slog.KindDuration:
		return attribute.Int64(key, int64(value.Duration()))
	case slog.KindFloat64:
		return attribute.Float64(key, value.Float64())
	case slog.KindInt64:
		return attribute.Int64(key, value.Int64())
	case slog.KindString:
		return attribute.String(key, value.String())
	case slog.KindTime:
		return attribute.String(key, value.Time().Format(time.RFC3339))
	case slog.KindGroup:
		return attribute.String(key, fmt.Sprintf("%v", value.Group()))
	default:
		return attribute.String(key, value.String())
	}
}

func durationToMillisecond(d time.Duration) float64 {
	return float64(d/1000) / 1000
}

func durationToSecond(d time.Duration) float64 {
	return float64(d/1000) / 1000000
}
