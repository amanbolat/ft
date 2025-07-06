# ft – function trace

[![Go Report Card](https://goreportcard.com/badge/github.com/amanbolat/ft)](https://goreportcard.com/report/github.com/amanbolat/ft)
[![GoDoc](https://godoc.org/github.com/amanbolat/ft?status.svg)](https://godoc.org/github.com/amanbolat/ft)
[![License](https://img.shields.io/badge/license-BSD%20Zero%20Clause%20License-blue.svg)](https://opensource.org/license/0bsd/)

A lightweight library for tracing function execution with OpenTelemetry integration,
structured logging, and metrics collection.

## Why?

In most of my projects, I start by relying on simple logs to get quick insights into the application. 
As the project evolves, I usually add metrics and traces for improved observability, 
but setting up all the libraries can be time-consuming. 
This library allows you to enable the observability of your functions with just two lines of code.

## Features

- Structured logging using `slog`.
- OpenTelemetry tracing integration.
- Metrics for execution counts and duration. 
- Configurable log level.
- Option to opt out of metrics and tracing.

## Usage

Just add two lines to the beginning of the function:

```go
func Do(ctx context.Context) (err error) {
    ctx, span := ft.Start(ctx, "main.Do", ft.WithErr(&err)) // Log when we enter the `Do` function.
    defer span.End()                                        // Log, trace and meter when we exit the `Do` function.

    err = errors.New("unexpected error")
	
    return
}
```

If you run the code above, you will see the following output:

```shell
time=2025-01-25T21:55:27.068+01:00 level=INFO msg="action started" action=main.Do
time=2025-01-25T21:55:27.069+01:00 level=ERROR msg="action ended" action=main.Do duration_ms=0.743 error="unexpected error"
```

### Adding attributes dynamically

You can add attributes to a span after it has been created using the `AddAttrs` method. This is useful when you want to add contextual information that becomes available during function execution:

```go
func ProcessUser(ctx context.Context, userID string) (err error) {
    ctx, span := ft.Start(ctx, "main.ProcessUser", ft.WithErr(&err))
    defer span.End()

    // Add initial attributes
    span.AddAttrs(slog.String("user_id", userID))

    user, err := fetchUser(userID)
    if err != nil {
        return err
    }

    // Add more attributes as they become available
    span.AddAttrs(
        slog.String("user_email", user.Email),
        slog.Bool("user_premium", user.IsPremium),
    )

    // Process user...
    return nil
}
```

The `AddAttrs` method is thread-safe and can be called multiple times throughout the function execution. All attributes will be included in the final log output and, if enabled, added to the OpenTelemetry span.

### OpenTelemetry Integration

Setup OTEL tracer and meter globally and `ft` will start sending metrics and traces to the OTLP collector:

```go
mp, _, _ := autometer.NewMeterProvider(context.Background())
otel.SetMeterProvider(mp)
tp, _, _ := autotracer.NewTracerProvider(context.Background())
otel.SetTracerProvider(tp)

ft.SetMetricsEnabled(true)
ft.SetTracingEnabled(true)
```

> [!TIP]
> In the example above I use [go-faster/sdk](https://github.com/go-faster/sdk) to setup OTEL based on environment
> variables.


## Configuration

`ft` package provides many functions to configure its behaviour. See the table below:

Here's a markdown table with all the `Set` functions and their descriptions:

| Function                                 | Description                                                                                                                                                  |
|------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `SetDurationMetricUnit(unit string)`     | Sets the global duration metric unit. Accepts either millisecond (`ms`) or second (`s`) as valid units. Defaults to millisecond if invalid unit is provided. |
| `SetDefaultLogger(l *slog.Logger)`       | Sets the global logger instance. Does nothing if nil logger is provided.                                                                                     |
| `SetLogLevelOnFailure(level slog.Level)` | Sets the global log level for failure scenarios.                                                                                                             |
| `SetLogLevelOnSuccess(level slog.Level)` | Sets the global log level for success scenarios.                                                                                                             |
| `SetTracingEnabled(v bool)`              | Enables or disables global tracing functionality.                                                                                                            |
| `SetMetricsEnabled(v bool)`              | Enables or disables global metrics collection.                                                                                                               |
| `SetClock(c clockwork.Clock)`            | Sets the global clock instance used for time-related operations.                                                                                             |
| `SetAppendOtelAttrs(v bool)`             | Enables or disables the appending of OpenTelemetry attributes globally.                                                                                      |
