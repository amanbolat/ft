package ft_test

import (
	"context"
	"errors"
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
	// Setup
	spanRecorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
	otel.SetTracerProvider(tp)
	
	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)
	
	ft.SetTracingEnabled(true)
	
	// Test
	ctx := context.Background()
	testAction := "test_action"
	
	ctx, span := ft.Start(ctx, testAction)
	
	// Advance time by 100ms
	fakeClock.Advance(100 * time.Millisecond)
	
	span.End()
	
	// Verify
	spans := spanRecorder.Ended()
	require.Len(t, spans, 1)
	
	recordedSpan := spans[0]
	assert.Equal(t, testAction, recordedSpan.Name())
	assert.Equal(t, codes.Unset, recordedSpan.Status().Code)
	
	// Verify duration
	assert.Equal(t, 100*time.Millisecond, recordedSpan.EndTime().Sub(recordedSpan.StartTime()))
}

func TestSpan_WithError(t *testing.T) {
	// Setup
	spanRecorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
	otel.SetTracerProvider(tp)
	
	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)
	
	ft.SetTracingEnabled(true)
	
	// Test
	ctx := context.Background()
	testAction := "test_action_error"
	testError := errors.New("test error")
	var err error = testError
	
	ctx, span := ft.Start(ctx, testAction, ft.WithErr(&err))
	
	fakeClock.Advance(50 * time.Millisecond)
	
	span.End()
	
	// Verify
	spans := spanRecorder.Ended()
	require.Len(t, spans, 1)
	
	recordedSpan := spans[0]
	assert.Equal(t, testAction, recordedSpan.Name())
	assert.Equal(t, codes.Error, recordedSpan.Status().Code)
	assert.Equal(t, testError.Error(), recordedSpan.Status().Description)
}

func TestSpan_WithAttributes(t *testing.T) {
	// Setup
	spanRecorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
	otel.SetTracerProvider(tp)
	
	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)
	
	ft.SetTracingEnabled(true)
	ft.SetAppendOtelAttrs(true)
	
	// Test
	ctx := context.Background()
	testAction := "test_action_attrs"
	
	ctx, span := ft.Start(ctx, testAction, ft.WithAttrs(
		slog.String("test_key", "test_value"),
		slog.Int64("test_number", 42),
	))
	
	span.End()
	
	// Verify
	spans := spanRecorder.Ended()
	require.Len(t, spans, 1)
	
	recordedSpan := spans[0]
	assert.Equal(t, testAction, recordedSpan.Name())
	
	expectedAttrs := []attribute.KeyValue{
		attribute.String("action", testAction),
		attribute.String("test_key", "test_value"),
		attribute.Int64("test_number", 42),
	}
	
	spanAttrs := recordedSpan.Attributes()
	assert.ElementsMatch(t, expectedAttrs, spanAttrs)
}

func TestSpan_TracingDisabled(t *testing.T) {
	// Setup
	spanRecorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
	otel.SetTracerProvider(tp)
	
	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)
	
	ft.SetTracingEnabled(false)
	
	// Test
	ctx := context.Background()
	testAction := "test_action_disabled"
	
	ctx, span := ft.Start(ctx, testAction)
	span.End()
	
	// Verify no spans were recorded
	spans := spanRecorder.Ended()
	assert.Empty(t, spans)
}

func TestSpan_MetricsEnabled(t *testing.T) {
	// Setup
	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)
	
	ft.SetMetricsEnabled(true)
	ft.SetDurationMetricUnit(ft.DurationMetricUnitMillisecond)
	
	mp := noop.NewMeterProvider()
	otel.SetMeterProvider(mp)
	
	// Test
	ctx := context.Background()
	testAction := "test_action_metrics"
	
	ctx, span := ft.Start(ctx, testAction)
	fakeClock.Advance(75 * time.Millisecond)
	span.End()
}

func TestSpan_NilContext(t *testing.T) {
	// Setup
	spanRecorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
	otel.SetTracerProvider(tp)
	
	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)
	
	ft.SetTracingEnabled(true)
	
	// Test with nil context
	testAction := "test_action_nil_ctx"
	
	ctx, span := ft.Start(nil, testAction)
	
	// Verify context was created
	assert.NotNil(t, ctx)
	
	span.End()
	
	// Verify span was recorded
	spans := spanRecorder.Ended()
	require.Len(t, spans, 1)
	assert.Equal(t, testAction, spans[0].Name())
}

func TestSpan_DurationUnits(t *testing.T) {
	// Setup
	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)
	
	ft.SetMetricsEnabled(true)
	
	// Test milliseconds
	ft.SetDurationMetricUnit(ft.DurationMetricUnitMillisecond)
	ctx, span := ft.Start(context.Background(), "test_ms")
	fakeClock.Advance(100 * time.Millisecond)
	span.End()
	
	// Test seconds
	ft.SetDurationMetricUnit(ft.DurationMetricUnitSecond)
	ctx, span = ft.Start(context.Background(), "test_s")
	fakeClock.Advance(1 * time.Second)
	span.End()
}

func TestSpan_LogLevels(t *testing.T) {
	// Setup
	fakeClock := clockwork.NewFakeClock()
	ft.SetClock(fakeClock)
	
	var logBuffer testLogBuffer
	logger := slog.New(slog.NewTextHandler(&logBuffer, nil))
	ft.SetDefaultLogger(logger)
	
	ft.SetLogLevelOnSuccess(slog.LevelDebug)
	ft.SetLogLevelOnFailure(slog.LevelError)
	
	// Test success case
	ctx, span := ft.Start(context.Background(), "test_success")
	span.End()
	assert.Contains(t, logBuffer.String(), "level=DEBUG")
	
	// Test failure case
	logBuffer.Reset()
	err := errors.New("test error")
	ctx, span = ft.Start(context.Background(), "test_failure", ft.WithErr(&err))
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
