package ft

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"

	"log/slog"
	"time"
)

func Trace(ctx context.Context, action string) Span {
	start := time.Now()

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

	if DefaultLogger().Enabled(context.Background(), slog.LevelInfo) {
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:])
		r := slog.NewRecord(start, logLevel.Level(), "action started", pcs[0])
		r.AddAttrs(attrs...)
		_ = DefaultLogger().Handler().Handle(context.Background(), r)

		DefaultLogger().LogAttrs(ctx, logLevel.Level(), "action started", attrs...)
	}

	return Span{
		ctx:       ctx,
		start:     start.UnixNano(),
		action:    action,
		traceSpan: span,
	}
}

type Span struct {
	ctx       context.Context
	start     int64
	action    string
	traceSpan trace.Span
}

func (s Span) Log() {
	duration := time.Now().UnixNano() - s.start

	DefaultLogger().LogAttrs(
		s.ctx, logLevel.Level(),
		"action ended",
		slog.String("action", s.action),
		slog.Duration("duration", time.Duration(duration)),
	)

	if s.traceSpan != nil {
		s.traceSpan.End()
	}
}
