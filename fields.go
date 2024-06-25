package errors

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
)

// HasFields Implement this interface to pass along unstructured context to the logger.
// It is the responsibility of Fields() implementation to unwrap the error chain and
// collect all errors that have `Fields()` defined.
type HasFields interface {
	Fields() []any
	Error() string
}

// HasFormat True if the interface has the format method (from fmt package)
type HasFormat interface {
	Format(st fmt.State, verb rune)
}

// Fields Creates errors that conform to the `Fields` interface
type Fields []any

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func (f Fields) Wrap(err error) error {
	if err == nil {
		return nil
	}
	return &fields{
		fields:  f,
		wrapped: err,
	}
}

func (f Fields) Error(msg string) error {
	return &fields{
		fields:  f,
		wrapped: errors.New(msg),
	}
}

func (f Fields) Errorf(format string, args ...any) error {
	return &fields{
		fields:  f,
		wrapped: fmt.Errorf(format, args...),
	}
}

type fields struct {
	fields  Fields
	wrapped error
}

func (c *fields) Unwrap() error {
	u, ok := c.wrapped.(interface {
		Unwrap() error
	})
	if !ok {
		return c.wrapped
	}
	return u.Unwrap()
}

func (c *fields) Is(target error) bool {
	_, ok := target.(*fields)
	return ok
}

func (c *fields) Error() string {
	return c.wrapped.Error()
}

func (c *fields) Fields() []any {
	var result []any
	result = append(result, c.fields...)

	// child fields have precedence as they are closer to the cause
	var f HasFields
	if errors.As(c.wrapped, &f) {
		child := f.Fields()
		if child == nil {
			return result
		}
		result = append(result, child...)
	}
	// child fields have precedence as they are closer to the cause
	return result
}

func (c *fields) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v (%s)", c.wrapped, c.FormatFields())
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, c.Error())
		return
	}
}

func (c *fields) FormatFields() string {
	var buf bytes.Buffer
	var count int

	args := c.Fields()
	var attr slog.Attr

	for len(args) > 0 {
		if count > 0 {
			buf.WriteString(", ")
		}
		attr, args = argsToAttr(args)
		buf.WriteString(fmt.Sprintf("%+v=%+v", attr.Key, attr.Value.Any()))
		count++
	}
	return buf.String()
}

// ToMap collects all the fields from any errors that may have been wrapped
func ToMap(err error) map[string]any {
	result := map[string]any{
		"err": err.Error(),
	}

	// Search the error chain for fields
	var f HasFields
	if errors.As(err, &f) {
		args := f.Fields()
		var attr slog.Attr

		for len(args) > 0 {
			attr, args = argsToAttr(args)
			result[attr.Key] = attr.Value.Any()
		}
	}
	return result
}

// ToAttr returns the field information for the underlying error as
// slog compatible arguments
//
//	err := errors.Fields{"key1", "value1"}.Error("error")
//	slog.Error("some error occurred", errors.ToAttr(err))
func ToAttr(err error) []any {
	result := []any{
		"err", err.Error(),
	}

	// Search the error chain for fields
	var f HasFields
	if errors.As(err, &f) {
		result = append(result, f.Fields()...)
	}
	return result
}

const badKey = "!BADKEY"

// argsToAttr turns a prefix of the nonempty args slice into an Attr
// and returns the unconsumed portion of the slice.
// If args[0] is an Attr, it returns it.
// If args[0] is a string, it treats the first two elements as
// a key-value pair.
// Otherwise, it treats args[0] as a value with a missing key.
func argsToAttr(args []any) (slog.Attr, []any) {
	switch x := args[0].(type) {
	case string:
		if len(args) == 1 {
			return slog.String(badKey, x), nil
		}
		return slog.Any(x, args[1]), args[2:]

	case slog.Attr:
		return x, args[1:]

	default:
		return slog.Any(badKey, x), args[1:]
	}
}
