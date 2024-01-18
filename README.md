# ft – log and trace a function

A very simple utility library written in Go to log and trace a function lifecycle.

## Why?

Because I'm tired of adding the same boilerplate code to every function I write, that usually requires three or more lines of code to log and trace a function and some metrics.

## Usage

Just add `defer ft.Trace(ctx, "<method name>").WithError(&err).Log()` to the beginning of the function.

```go
func Do(ctx context.Context) (err error) {
	defer ft.Trace(ctx, "main.Do").WithError(&err).Log()
	err = errors.New("unexpected error")

	return
}
```

If you run the code above, you will see the following output:

```shell
2024/01/18 01:05:07 INFO action started action=main.Do
2024/01/18 01:05:07 ERROR action ended action=main.Do duration=472.583µs error=unexpected error
```

## Tracing

Add the following to the main function to enable tracing:

```go
func main() {
    ft.EnableTracing()
}
```

## Metrics

Add the following to the main function to enable metrics:

```go
func main() {
    ft.EnableMetrics()
}
```

It uses [go-metrics](github.com/hashicorp/go-metrics) package under the hood.


## License

BSD Zero Clause License

Copyright (c) [2024] [Amanbolat Balabekov]

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
PERFORMANCE OF THIS SOFTWARE.
