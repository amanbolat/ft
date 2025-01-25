# ft – function trace

[![Go Report Card](https://goreportcard.com/badge/github.com/amanbolat/ft)](https://goreportcard.com/report/github.com/amanbolat/ft)
[![GoDoc](https://godoc.org/github.com/amanbolat/ft?status.svg)](https://godoc.org/github.com/amanbolat/ft)
[![License](https://img.shields.io/badge/license-BSD%20Zero%20Clause%20License-blue.svg)](https://opensource.org/license/0bsd/)

A very simple utility library written in Go to log and trace a function lifecycle.

## Why?

Because I'm tired of adding the same boilerplate code to every function I write, that usually requires three or more 
lines of code to log and trace a function and some metrics.

## Usage

Just add these two lines to the beginning of the function:

```go
ctx, span := ft.Start(ctx, "package.Function", ft.WithErr(&err))
defer span.End()
```

```go
func Do(ctx context.Context) (err error) {
    ctx, span := ft.Start(ctx, "main.Do", ft.WithErr(&err))
    defer span.End()
	
    err := callAnotherFn()
    if err != nil {
        return err
    }   

    return nil
}
```

If you run the code above, you will see the following output:

```shell
time=2025-01-20T00:32:20.025+01:00 level=INFO msg="action started" action=main.Do
time=2025-01-20T00:32:20.330+01:00 level=ERROR msg="action ended" action=main.Do duration_ms=137.546 error="error from callAnotherFn"
```


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
