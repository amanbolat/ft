package main

import (
	"context"
	"errors"

	"github.com/amanbolat/ft"
)

// Running this example will produce an output similar to this:
//
// time=2025-01-25T21:55:27.068+01:00 level=INFO msg="action started" action=main.Do
// time=2025-01-25T21:55:27.069+01:00 level=ERROR msg="action ended" action=main.Do duration_ms=0.743 error="unexpected error"
func main() {
	ctx := context.Background()
	_ = Do(ctx)
}

func Do(ctx context.Context) (err error) {
	ctx, span := ft.Start(ctx, "main.Do", ft.WithErr(&err))
	defer span.End()

	err = errors.New("unexpected error")

	return
}
