package ft

import (
	"context"
	"runtime"

	"github.com/hashicorp/go-metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"

	"log/slog"
	"time"
)

func Trace(ctx context.Context, action string) Span {
	now := time.Now()

	if ctx == nil {
		ctx = context.Background()
	}
	var span trace.Span

	if tracingEnabled.Load() {
		ctx, span = otel.Tracer(
			instrumentationName,
			trace.WithInstrumentationVersion("git_ref:"+gitCommit),
			trace.WithSchemaURL(semconv.SchemaURL),
		).Start(
			ctx,
			action,
			trace.WithSpanKind(trace.SpanKindInternal),
			trace.WithAttributes(
				attribute.String("action", action),
			),
		)
	}

	attrs := []slog.Attr{
		slog.String("action", action),
	}

	log(ctx, "action started", logLevel.Level(), now, attrs...)

	return Span{
		ctx:       ctx,
		start:     now,
		action:    action,
		traceSpan: span,
	}
}

type Span struct {
	ctx       context.Context
	start     time.Time
	action    string
	traceSpan trace.Span
	err       *error
}

func (s Span) WithError(err *error) Span {
	s.err = err

	return s
}

func (s Span) Log() {
	now := time.Now()
	level := logLevel.Level()

	attrs := []slog.Attr{
		slog.String("action", s.action),
		slog.Duration("duration", now.Sub(s.start)),
	}

	if s.err != nil && *s.err != nil {
		level = slog.LevelError
		attrs = append(attrs, slog.Any("error", *s.err))

		if s.traceSpan != nil {
			s.traceSpan.RecordError(*s.err)
		}
	}

	if metricsEnabled.Load() {
		metrics.MeasureSinceWithLabels([]string{s.action}, s.start, nil)
	}

	log(s.ctx, "action ended", level, now, attrs...)

	if s.traceSpan != nil {
		s.traceSpan.End()
	}
}

func log(ctx context.Context, msg string, level slog.Level, now time.Time, attrs ...slog.Attr) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(now, level, msg, pcs[0])
	r.AddAttrs(attrs...)
	_ = DefaultLogger().Handler().Handle(ctx, r)
}
