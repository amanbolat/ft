package ft_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/amanbolat/ft"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric/noop"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestSpan_Basic(t *testing.T) {
	ft.SetDefaultLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))
	spanRecorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
	otel.SetTracerProvider(tp)

	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)

	ft.SetTracingEnabled(true)

	ctx := context.Background()
	testAction := "test_action"

	ctx, span := ft.Start(ctx, testAction)
	fakeClock.Advance(100 * time.Millisecond)
	span.End()

	spans := spanRecorder.Ended()
	require.Len(t, spans, 1)

	recordedSpan := spans[0]
	assert.Equal(t, testAction, recordedSpan.Name())
	assert.Equal(t, codes.Unset, recordedSpan.Status().Code)
	assert.Equal(t, 100*time.Millisecond, recordedSpan.EndTime().Sub(recordedSpan.StartTime()))
}

func TestSpan_WithError(t *testing.T) {
	ft.SetDefaultLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))
	spanRecorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
	otel.SetTracerProvider(tp)

	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)
	ft.SetTracingEnabled(true)

	ctx := context.Background()
	testAction := "test_action_error"
	testError := errors.New("test error")
	var err error = testError

	ctx, span := ft.Start(ctx, testAction, ft.WithErr(&err))
	fakeClock.Advance(50 * time.Millisecond)
	span.End()

	spans := spanRecorder.Ended()
	require.Len(t, spans, 1)

	recordedSpan := spans[0]
	assert.Equal(t, testAction, recordedSpan.Name())
	assert.Equal(t, codes.Error, recordedSpan.Status().Code)
	assert.Equal(t, testError.Error(), recordedSpan.Status().Description)
}

type CustomValue struct {
	val1 string
	val2 string
}

func (v CustomValue) LogValue() slog.Value {
	return slog.StringValue(v.val1 + "_" + v.val2)
}

func TestSpan_WithAttributes(t *testing.T) {
	ft.SetDefaultLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))
	spanRecorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
	otel.SetTracerProvider(tp)

	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)
	ft.SetTracingEnabled(true)
	ft.SetAppendOtelAttrs(true)

	ctx := context.Background()
	testAction := "test_action_attrs"

	ctx, span := ft.Start(ctx, testAction, ft.WithAttrs(
		slog.String("string", "value"),
		slog.Int64("int", 1),
		slog.Bool("bool", true),
		slog.Duration("duration", time.Second*500),
		slog.Float64("float", 2),
		slog.Time("timestamp", time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC)),
		slog.Group("group", slog.String("string", "value"), slog.Int64("int", 1)),
		slog.Any("custom_value", CustomValue{val1: "a", val2: "b"}),
	))

	span.End()

	spans := spanRecorder.Ended()
	require.Len(t, spans, 1)

	recordedSpan := spans[0]
	assert.Equal(t, testAction, recordedSpan.Name())

	expectedAttrs := []attribute.KeyValue{
		attribute.String("action", testAction),
		attribute.String("string", "value"),
		attribute.Int64("int", 1),
		attribute.Bool("bool", true),
		attribute.Int64("duration", int64(time.Second*500/1000/1000)),
		attribute.Float64("float", 2),
		attribute.String("timestamp", time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC).Format(time.RFC3339)),
		attribute.String("group", "[string=value int=1]"),
		attribute.String("custom_value", "a_b"),
	}

	spanAttrs := recordedSpan.Attributes()
	assert.ElementsMatch(t, expectedAttrs, spanAttrs)
}

func TestSpan_TracingDisabled(t *testing.T) {
	ft.SetDefaultLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))
	spanRecorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
	otel.SetTracerProvider(tp)

	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)
	ft.SetTracingEnabled(false)

	ctx := context.Background()
	testAction := "test_action_disabled"

	ctx, span := ft.Start(ctx, testAction)
	span.End()

	spans := spanRecorder.Ended()
	assert.Empty(t, spans)
}

func TestSpan_MetricsEnabled(t *testing.T) {
	ft.SetDefaultLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))
	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)

	ft.SetMetricsEnabled(true)
	ft.SetDurationMetricUnit(ft.DurationMetricUnitMillisecond)

	mp := noop.NewMeterProvider()
	otel.SetMeterProvider(mp)

	ctx := context.Background()
	testAction := "test_action_metrics"

	ctx, span := ft.Start(ctx, testAction)
	fakeClock.Advance(75 * time.Millisecond)
	span.End()
}

func TestSpan_NilContext(t *testing.T) {
	ft.SetDefaultLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))
	spanRecorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
	otel.SetTracerProvider(tp)

	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)
	ft.SetTracingEnabled(true)

	testAction := "test_action_nil_ctx"
	ctx, span := ft.Start(nil, testAction)
	assert.NotNil(t, ctx)

	span.End()

	spans := spanRecorder.Ended()
	require.Len(t, spans, 1)
	assert.Equal(t, testAction, spans[0].Name())
}

func TestSpan_DurationUnits(t *testing.T) {
	ft.SetDefaultLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))
	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)

	ft.SetMetricsEnabled(true)

	// Test with milliseconds
	ft.SetDurationMetricUnit(ft.DurationMetricUnitMillisecond)
	var logBuffer testLogBuffer
	logger := slog.New(slog.NewTextHandler(&logBuffer, nil))
	ft.SetDefaultLogger(logger)
	
	_, span := ft.Start(context.Background(), "test_ms")
	fakeClock.Advance(100 * time.Millisecond)
	span.End()
	
	// Verify millisecond format in logs
	assert.Contains(t, logBuffer.String(), "duration_ms=")
	
	// Test with seconds
	logBuffer.Reset()
	ft.SetDurationMetricUnit(ft.DurationMetricUnitSecond)
	_, span = ft.Start(context.Background(), "test_s")
	fakeClock.Advance(1 * time.Second)
	span.End()
	
	// Verify second format in logs
	assert.Contains(t, logBuffer.String(), "duration_s=")
}

func TestSpan_LogLevels(t *testing.T) {
	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)

	var logBuffer testLogBuffer
	logger := slog.New(slog.NewTextHandler(&logBuffer, nil))
	ft.SetDefaultLogger(logger)

	ft.SetLogLevelOnSuccess(slog.LevelDebug)
	ft.SetLogLevelOnFailure(slog.LevelError)

	_, span := ft.Start(context.Background(), "test_success")
	span.End()
	assert.Contains(t, logBuffer.String(), "level=DEBUG")

	logBuffer.Reset()
	err := errors.New("test error")
	_, span = ft.Start(context.Background(), "test_failure", ft.WithErr(&err))
	span.End()
	assert.Contains(t, logBuffer.String(), "level=ERROR")
}

type testLogBuffer struct {
	content string
}

func (b *testLogBuffer) Write(p []byte) (n int, err error) {
	b.content += string(p)
	return len(p), nil
}

func (b *testLogBuffer) String() string {
	return b.content
}

func (b *testLogBuffer) Reset() {
	b.content = ""
}
