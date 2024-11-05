package errors_test

import (
	"context"
	"fmt"
	"github.com/kapetan-io/errors"
	"log/slog"
)

type Struct struct {
	attrs *errors.Attrs
}

func ExampleAttrs() {
	// Attach context to an error
	err := errors.With("foo", "bar").Error("query failed")

	// Prints `query failed`
	fmt.Printf("%s\n", err)

	// Prints `query failed (foo=bar)`
	fmt.Printf("%+v\n", err)

	// Prints `Attributes [foo=bar]`
	fmt.Printf("Attributes %v\n", errors.AttrsFrom(err))

	// Prints `2024/09/30 12:00:00 level=ERROR msg="query failed" foo=bar`
	slog.LogAttrs(context.Background(), slog.LevelError, err.Error(), errors.AttrsFrom(err)...)

	// Use with slog.Attr to create complex logging context
	s := Struct{
		attrs: errors.WithAttr(slog.String("database", "sqlite")),
	}

	// Prints `query failed (database=sqlite, foo=bar)`
	fmt.Printf("%+v\n", s.attrs.With("foo", "bar").Error("query failed"))

	// Which can be used to attach to logging
	slog.LogAttrs(context.Background(), slog.LevelError, err.Error(), errors.AttrsFromWithCodeLoc(err)...)

	// Prints `message: query error (key1=value1)`
	err = errors.New("query error")
	wrap := errors.With("key1", "value1").Errorf("message: %w", err)
	fmt.Printf("%+v\n", wrap)

	// Output:
	// query failed
	// query failed (foo=bar)
	// Attributes [foo=bar]
	// query failed (database=sqlite, foo=bar)
	// message: query error (key1=value1)
}
