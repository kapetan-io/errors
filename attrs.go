package errors

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"runtime"
)

// HasAttrs is used identify which errors have attributes attached in order to pass along unstructured
// context down the err tree. It is the responsibility of Attrs() implementation to unwrap the error
// tree and collect all errors that implement this interface.
type HasAttrs interface {
	Attrs() ([]slog.Attr, uintptr)
	Error() string
}

// With returns an *Attrs which includes the given attributes
func With(args ...any) *Attrs {
	slog.With()
	a := &Attrs{}
	return a.With(args...)
}

// WithAttr returns an *Attrs which includes the given slog.Attr
func WithAttr(attrs ...slog.Attr) *Attrs {
	a := &Attrs{}
	return a.WithAttr(attrs...)
}

// Error works exactly like standard lib `errors.New()` and includes
// stack information which can be extracted with errors.AttrsWithCodeLoc()
// or ErrAttrs.Attrs()
func Error(msg string) error {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [runtime.Callers, and this function]
	return &ErrAttrs{
		wrapped: errors.New(msg),
		attrs:   &Attrs{},
		pc:      pcs[0],
	}
}

// Errorf works exactly like standard lib `fmt.Errorf()` and includes
// stack information which can be extracted with errors.AttrsWithCodeLoc()
// or ErrAttrs.Attrs()
func Errorf(format string, args ...any) error {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [runtime.Callers, and this function]
	return &ErrAttrs{
		wrapped: fmt.Errorf(format, args...),
		attrs:   &Attrs{},
		pc:      pcs[0],
	}
}

// Wrap returns an error with stack information for the code location where Wrap
// is called. The returned error has no "message" but defers to the wrapped message
// when Error() is called. If err is nil, Wrap returns nil. Stack information can
// be extracted with errors.AttrsWithCodeLoc() or ErrAttrs.Attrs()
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [runtime.Callers, and this function]
	return &ErrAttrs{
		attrs:   &Attrs{},
		pc:      pcs[0],
		wrapped: err,
	}
}

// Logger is a crazy idea which would extract the attributes from
// the currently configured logger.
//
//  log := slog.With("foo", "bar")
//  return errors.Logger(log).Error("my error message")
//
//func Logger(l *slog.Logger) *Attrs {
//  a := &Attrs{}
//	// Get the attributes somehow?
//	return a.WithAttr(l.GetAttr())
//}

// Attrs holds attached attributes until Error() or Errorf() are called to
// return the attributes via ErrAttrs as an error.
type Attrs struct {
	attrs []slog.Attr
}

// With returns a new *Attrs which includes the given attributes combined
// with any existing attributes defined in the current Attrs.
func (a *Attrs) With(args ...any) *Attrs {
	return a.WithAttr(argsToAttrSlice(args)...)
}

// WithAttr returns a new *Attrs which includes the given attributes combined
// with any existing attributes defined in the current Attrs.
func (a *Attrs) WithAttr(as ...slog.Attr) *Attrs {
	return &Attrs{attrs: append(a.attrs, as...)}
}

// Wrap returns an error with included code location information
// at the point Wrap is called. The returned error has no "message" but
// defers to the wrapped message when Error() is called.
// If err is nil, Wrap returns nil.
func (a *Attrs) Wrap(err error) error {
	if err == nil {
		return nil
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [runtime.Callers, and this function]
	return &ErrAttrs{
		pc:      pcs[0],
		wrapped: err,
		attrs:   a,
	}
}

// Error works exactly like standard lib `errors.New()` and includes
// stack information which can be extracted with errors.AttrsWithCodeLoc()
// or ErrAttrs.Attrs()
func (a *Attrs) Error(msg string) error {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [runtime.Callers, and this function]
	return &ErrAttrs{
		wrapped: errors.New(msg),
		pc:      pcs[0],
		attrs:   a,
	}
}

// Errorf works exactly like standard lib `fmt.Errorf()` and includes
// stack information which can be extracted with errors.AttrsWithCodeLoc()
// or ErrAttrs.Attrs()
func (a *Attrs) Errorf(format string, args ...any) error {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [runtime.Callers, and this function]
	return &ErrAttrs{
		wrapped: fmt.Errorf(format, args...),
		pc:      pcs[0],
		attrs:   a,
	}
}

// ErrAttrs is an error which has slog.Attr attached
type ErrAttrs struct {
	pc      uintptr
	attrs   *Attrs
	wrapped error
}

// Error returns the error as a string
func (e *ErrAttrs) Error() string {
	return e.wrapped.Error()
}

// Is returns true if the target is of type ErrAttrs
func (e *ErrAttrs) Is(target error) bool {
	_, ok := target.(*ErrAttrs)
	return ok
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func (e *ErrAttrs) Unwrap() error {
	u, ok := e.wrapped.(interface {
		Unwrap() error
	})
	if !ok {
		return e.wrapped
	}
	return u.Unwrap()
}

// Attrs recursively returns all attributes in the err tree.
// The pc returned is from the ErrAttrs closest to the root of the
// err tree.
func (e *ErrAttrs) Attrs() ([]slog.Attr, uintptr) {
	var result []slog.Attr
	result = append(result, e.attrs.attrs...)
	pc := e.pc

	var (
		child []slog.Attr
		a     HasAttrs
	)
	if errors.As(e.wrapped, &a) {
		child, pc = a.Attrs()
		if child == nil {
			return result, pc
		}
		result = append(result, child...)
	}
	return result, pc
}

// Format follows the standard set forth by the fmt package
// for serializing structures using formating directives %s, %v, %+v, %q
func (e *ErrAttrs) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v (%s)", e.wrapped, e.formatAttrs())
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
		return
	}
}

func (e *ErrAttrs) formatAttrs() string {
	var buf bytes.Buffer
	var count int

	attrs, _ := e.Attrs()
	for _, attr := range attrs {
		if count > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%+v=%+v", attr.Key, attr.Value.Any()))
		count++
	}
	return buf.String()
}

// AttrsFrom returns any attrs from the err tree. If the err tree contains
// no instances of HasAttrs then []slog.Attr{slog.Any("", nil)} is returned.
// This means it is safe to call with `slog.LogAttrs()` even if there are no
// attributes in the err tree.
func AttrsFrom(err error) []slog.Attr {
	var a HasAttrs
	if errors.As(err, &a) {
		attrs, _ := a.Attrs()
		return attrs
	}
	return []slog.Attr{slog.Any("", nil)}
}

// AttrsWithCodeLoc returns any attrs from the err tree and includes source code from the
// code position where the ErrAttrs error was created. The following OTEL fields
// are included in the returned slog.Attr returned.
//
//	code.filepath //path/to/file.go
//	code.function Struct.Method
//	code.lineno 156
//
// If the err tree contains no instances of HasAttrs then
// []slog.Attr{slog.Any("", nil)} is returned.
func AttrsWithCodeLoc(err error) []slog.Attr {
	var a HasAttrs
	if errors.As(err, &a) {
		attrs, pc := a.Attrs()
		attrs = append(attrs, attrsFromPC(pc)...)
		return attrs
	}
	return []slog.Attr{slog.Any("", nil)}
}

// --------------------------
// Private methods
// --------------------------

func attrsFromPC(pc uintptr) []slog.Attr {
	f, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	return []slog.Attr{
		slog.String(OtelCodeFilePath, f.File),
		slog.String(OtelCodeFunction, f.Function),
		slog.Int(OtelCodeLineNo, f.Line),
	}
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

func argsToAttrSlice(args []any) []slog.Attr {
	var (
		attr  slog.Attr
		attrs []slog.Attr
	)
	for len(args) > 0 {
		attr, args = argsToAttr(args)
		attrs = append(attrs, attr)
	}
	return attrs
}
