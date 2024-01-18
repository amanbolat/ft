package main

import (
	"context"
	"errors"

	"github.com/amanbolat/ft"
)

// Running this example will produce an output similar to this:
//
// 2024/01/18 01:05:07 INFO action started action=main.Do
// 2024/01/18 01:05:07 ERROR action ended action=main.Do duration=472.583µs error=unexpected error
func main() {
	ctx := context.Background()
	_ = Do(ctx)

}

func Do(ctx context.Context) (err error) {
	defer ft.Trace(ctx, "main.Do").WithError(&err).Log()

	err = errors.New("unexpected error")

	return
}
