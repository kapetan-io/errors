# Errors
Is a drop in replacement library for std "errors" package with support for optional fields 
to support the [only handle errors once rule](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
while not losing context where the error occurred.

## Usage
Attach fields to an error message using standard types or optionally use `slog.Attr` attributes.
```go

err := errors.Fields{
    "key1", "value1",
    slog.String("key2", "value2"),
    slog.Duration("expired", time.Second),
}.Errorf("message")

// Extract those attributes when reporting errors via slog
slog.Info("an error occurred!", errors.ToAttr(err))
// INFO an error occurred! err=message key1=value1 key2=value2 expired=1s
```

Use standard introspection functions to extract fields
```go
var f errors.HasFields
if errors.As(wrap, &f) {
    slog.Error("this error has log fields", f.Fields()...)
}
```

Each call on the call stack can add to the list of fields as needed.

## Convenience to std error library methods
Provides pass through access to the standard `errors.Is()`, `errors.As()`, `errors.Unwrap()` so you don't need to
import this package and the standard error package.


## Support for standard golang introspection functions
Errors wrapped with `errors.Fields{}` are compatible with standard library introspection functions `errors.Unwrap()`,
`errors.Is()` and `errors.As()`
```go
ErrQuery := errors.New("query error")
wrap := errors.Fields{"key1", "value1"}.Errorf("message: %w", err)
errors.Is(wrap, ErrQuery) // == true
```

## Custom Error types
You can even implement custom error types which can pass along fields by implementing the `errors.HasFields` interface.
```go
type ErrHasFields struct {
	M string
	F []any
}

func (e *ErrHasFields) Error() string {
	return e.M
}

func (e *ErrHasFields) Fields() []any {
	return e.F
}
```