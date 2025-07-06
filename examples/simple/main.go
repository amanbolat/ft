package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/amanbolat/ft"
)

// Running this example will produce an output similar to this:
//
// time=2025-01-25T21:55:27.068+01:00 level=INFO msg="action started" action=main.Do user_id=12345
// time=2025-01-25T21:55:27.069+01:00 level=ERROR msg="action ended" action=main.Do user_id=12345 processing_step=validation duration_ms=0.743 error="unexpected error"
func main() {
	ctx := context.Background()
	_ = Do(ctx)
}

func Do(ctx context.Context) (err error) {
	ctx, span := ft.Start(ctx, "main.Do", ft.WithErr(&err))
	defer span.End()

	// Add initial context
	span.AddAttrs(slog.String("user_id", "12345"))

	// Simulate some processing steps
	if err = validateInput(); err != nil {
		// Add information about where the error occurred
		span.AddAttrs(slog.String("processing_step", "validation"))
		return err
	}

	return nil
}

func validateInput() error {
	return errors.New("unexpected error")
}
